package basecamp3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://launchpad.37signals.com/authorization/new",
	TokenURL: "https://launchpad.37signals.com/authorization/token",
}

const (
	BasecampApiRootURL = "https://3.basecampapi.com"
)

type Basecamp struct {
	client          *http.Client
	oauth           *oauth2.Config
	TokenPersitence func(w http.ResponseWriter, r *http.Request) ContextWithTokenPersistence
}

type ContextWithTokenPersistence interface {
	context.Context
	oauth2.TokenSource
	SetToken(*oauth2.Token)
}

func (bc *Basecamp) getClient(tokensource oauth2.TokenSource) *http.Client {
	if bc.client == nil {
		bc.client = oauth2.NewClient(context.Background(), tokensource)
	}
	return bc.client
}

func (bc *Basecamp) verified(ctx ContextWithTokenPersistence, code string) error {
	tok, err := bc.oauth.Exchange(ctx, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		return err
	}
	ctx.SetToken(tok)
	bc.getClient(bc.oauth.TokenSource(ctx, tok))
	return nil
}

func (bc *Basecamp) VerifyRequest(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	err := bc.verified(bc.TokenPersitence(w, r), code)
	if err == nil {
		fmt.Fprintf(w, "Auth successful")
	} else {
		fmt.Fprintf(w, "Auth error: %v", err)
	}
}

func (bc *Basecamp) Authenticate(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, bc.oauth.AuthCodeURL("state", oauth2.SetAuthURLParam("type", "web_server")), http.StatusTemporaryRedirect)
}

func (bc *Basecamp) ApiReverseProxy(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/proxy")
	url := BasecampApiRootURL + path
	log.Printf("GET %s\n", url)
	bc.proxyGet(bc.TokenPersitence(w, r), url, w)
}

func extractToken(client *http.Client) *oauth2.Token {
	t, err := client.Transport.(*oauth2.Transport).Source.Token()
	if err == nil {
		return t
	}
	return nil
}

func (bc *Basecamp) get(ctx ContextWithTokenPersistence, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Balazsgrill's BC3 integration (https://github.com/balazsgrill/basecamp3)")
	client := bc.getClient(ctx)
	t := extractToken(client)
	if t == nil {
		// TODO redirect to auth
		return nil, errors.New("not authenticated")
	} else {
		ctx.SetToken(t)
		return client.Do(req)
	}
}

func (bc *Basecamp) getWithRateLimit(ctx ContextWithTokenPersistence, url string) (*http.Response, error) {
	resp, err := bc.get(ctx, url)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == 429 {
		// API Rate limit
		waits := resp.Header.Get("Retry-After")
		wait, err := strconv.Atoi(waits)
		if err != nil {
			log.Printf("HTTP 429 received, but no walid Retry-After header was provided (%s)", waits)
		}
		time.Sleep(time.Duration(wait) * time.Second)
		resp, err = bc.get(ctx, url)
		if err != nil {
			return resp, err
		}
	}
	return resp, err
}

func (bc *Basecamp) jsonGet(ctx ContextWithTokenPersistence, url string, value interface{}) error {
	resp, err := bc.getWithRateLimit(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return json.Unmarshal(data, value)
}

func (bc *Basecamp) proxyGet(ctx ContextWithTokenPersistence, url string, w http.ResponseWriter) {
	resp, err := bc.get(ctx, url)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func New(clientID string, clientSecret string, redirectURL string) *Basecamp {
	return &Basecamp{
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{},
			Endpoint:     Endpoint,
		},
		TokenPersitence: ReverseProxyRequestContext,
	}
}

// ServeLocalhost configures and starts a HTTP server on localhost which can be used to authenticate and query BC.
// Capable of handling a single BC session, and assumes the redirect URL to be set as http://localhost:<port>/verify
func ServeLocalhost(port int, clientID string, clientSecret string) {
	bc := New(clientID, clientSecret, fmt.Sprintf("http://localhost:%d/verify", port))

	// handle route using handler function
	http.HandleFunc("/verify", bc.VerifyRequest)
	http.HandleFunc("/auth", bc.Authenticate)
	http.HandleFunc("/proxy/", bc.ApiReverseProxy)

	log.Printf("Listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
