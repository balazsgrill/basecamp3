package basecamp3

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	client *http.Client
	oauth  *oauth2.Config
}

func (bc *Basecamp) verified(code string) error {
	ctx := context.Background()
	tok, err := bc.oauth.Exchange(ctx, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		return err
	}
	bc.client = bc.oauth.Client(ctx, tok)
	return nil
}

func (bc *Basecamp) VerifyRequest(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	err := bc.verified(code)
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
	path := strings.TrimPrefix(r.URL.Path, "/import")
	url := BasecampApiRootURL + path
	log.Printf("GET %s\n", url)
	bc.proxyGet(url, w)
}

func (bc *Basecamp) todos(project int, todolist int) (url string) {
	return fmt.Sprintf("%s/buckets/%d/todolists/%d/todos.json", BasecampApiRootURL, project, todolist)
}

func (bc *Basecamp) proxyGet(url string, w http.ResponseWriter) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, err)
		return
	}
	req.Header.Set("User-Agent", "Balazsgrill's BC3 integration (https://github.com/balazsgrill/basecamp)")
	resp, err := bc.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	return
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
