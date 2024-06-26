package logto

import (
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"time"
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

type Claims struct {
	JwtID    string    `json:"jti"`
	Subject  string    `json:"sub"`
	IssuedAt time.Time `json:"iat"`
	Exp      time.Time `json:"exp"`
	Scope    string    `json:"scope"`
	ClientId string    `json:"client_id"`
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
}

func (l *Logto) Claims(token *jwt.Token) *Claims {
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return l.defaultClaims(token)
	}
	claims := &Claims{
		JwtID:    cast.ToString(mapClaims["jti"]),
		Subject:  cast.ToString(mapClaims["sub"]),
		IssuedAt: time.Time{},
		Exp:      time.Time{},
		ClientId: cast.ToString(mapClaims["client_id"]),
		Issuer:   cast.ToString(mapClaims["iss"]),
		Audience: cast.ToString(mapClaims["aud"]),
	}
	var (
		exp interface{}
		iat interface{}
	)
	if exp, ok = mapClaims["exp"]; ok {
		claims.Exp = time.Unix(cast.ToInt64(exp), 0)
	}
	if iat, ok = mapClaims["iat"]; ok {
		claims.IssuedAt = time.Unix(cast.ToInt64(iat), 0)
	}
	return claims
}

func (l *Logto) defaultClaims(token *jwt.Token) *Claims {
	claims := &Claims{}
	sub, _ := token.Claims.GetSubject()
	claims.Subject = sub

	iat, _ := token.Claims.GetIssuedAt()
	if iat != nil {
		claims.IssuedAt = iat.Time
	}

	exp, _ := token.Claims.GetExpirationTime()
	if exp != nil {
		claims.Exp = exp.Time
	}

	iss, _ := token.Claims.GetIssuer()
	claims.Issuer = iss

	aud, _ := token.Claims.GetAudience()
	claims.Audience = aud[0]
	return claims
}
