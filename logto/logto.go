package logto

import (
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
)

type Logto struct {
	endpoint string
	keyFunc  keyfunc.Keyfunc
}

func NewLogto(endpoint string) (*Logto, error) {
	client := &Logto{
		endpoint: endpoint,
	}
	if err := client.init(endpoint); err != nil {
		return nil, err
	}
	return client, nil
}

func (l *Logto) init(endpoint string) error {
	jwksURI := fmt.Sprintf("https://%s/oidc/jwks", endpoint)

	k, err := keyfunc.NewDefault([]string{jwksURI}) // Context is used to end the refresh goroutine.
	if err == nil {
		l.keyFunc = k
		return nil
	}

	client := http.Client{}
	resp, err := client.Get(jwksURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	jwksJSON := json.RawMessage(b)

	k, err = keyfunc.NewJWKSetJSON(jwksJSON)
	if err != nil {
		return err
	}
	l.keyFunc = k
	return nil
}

func (l *Logto) Parse(token string) (*jwt.Token, error) {
	if l.keyFunc == nil {
		return nil, fmt.Errorf("keyFunc is nil")
	}
	return jwt.Parse(token, l.keyFunc.Keyfunc)
}
