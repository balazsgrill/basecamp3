package basecamp3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

const TokenCookie = "BC3-API-Token"

func getSavedToken(request *http.Request) *oauth2.Token {
	c, err := request.Cookie(TokenCookie)
	if err != nil {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return nil
	}
	t := &oauth2.Token{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return nil
	}
	return t
}

func saveToken(w http.ResponseWriter, t *oauth2.Token) {
	if t == nil {
		// TODO delete cookie
		return
	}
	data, _ := json.Marshal(t)
	tokenstr := base64.StdEncoding.EncodeToString(data)
	c := &http.Cookie{
		Name:    TokenCookie,
		Expires: t.Expiry,
		Value:   tokenstr,
	}
	http.SetCookie(w, c)
}

type cookieTokenContext struct {
	context.Context
	request *http.Request
	w       http.ResponseWriter
}

func ReverseProxyRequestContext(w http.ResponseWriter, r *http.Request) ContextWithTokenPersistence {
	return TokenStoredInCookieContext(context.Background(), w, r)
}

func TokenStoredInCookieContext(ctx context.Context, w http.ResponseWriter, r *http.Request) ContextWithTokenPersistence {
	return &cookieTokenContext{
		Context: ctx,
		request: r,
		w:       w,
	}
}

func (c *cookieTokenContext) Token() (*oauth2.Token, error) {
	return getSavedToken(c.request), nil
}

func (c *cookieTokenContext) SetToken(t *oauth2.Token) {
	saveToken(c.w, t)
}
