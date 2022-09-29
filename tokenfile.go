package basecamp3

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type tokenFile struct {
	context.Context
	filename string
}

func FileBasedTokenPersistence(filename string) func(w http.ResponseWriter, r *http.Request) ContextWithTokenPersistence {
	return func(w http.ResponseWriter, r *http.Request) ContextWithTokenPersistence {
		return TokenFile(context.Background(), filename)
	}
}

func TokenFile(ctx context.Context, filename string) ContextWithTokenPersistence {
	return &tokenFile{
		Context:  ctx,
		filename: filename,
	}
}

func (c *tokenFile) Token() (*oauth2.Token, error) {
	data, err := os.ReadFile(c.filename)
	if err == os.ErrNotExist {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *tokenFile) SetToken(t *oauth2.Token) {
	if t == nil {
		os.Remove(c.filename)
	} else {
		data, err := json.Marshal(t)
		if err == nil {
			err = os.WriteFile(c.filename, data, 0666)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}
}
