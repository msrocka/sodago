package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

// tokenTTL is the time a generated token is valid.
const tokenTTL = 90 * 24 * time.Hour

var (
	errInvalidToken = errors.New("invalid token")
	errExpiredToken = errors.New("expired token")
)

// tokenAuth generates and verifies the authentication tokens of the server. A
// token is a JSON web token that is signed with a key stored in the data
// folder, so that tokens stay valid when the server restarts.
type tokenAuth struct {
	key []byte
}

// tokenClaims are the claims that are stored in a token.
type tokenClaims struct {
	Subject string `json:"sub"`
	Issued  int64  `json:"iat"`
	Expires int64  `json:"exp"`
}

// initTokenAuth creates the token handler and ensures that the signing key
// exists in the data folder.
func initTokenAuth(args Args) (*tokenAuth, error) {
	path := filepath.Join(args.DataDir(), "token_auth.key")
	key, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		log.Println("Init token auth ...")
		key = securecookie.GenerateRandomKey(32)
		if err := os.WriteFile(path, key, os.ModePerm); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &tokenAuth{key: key}, nil
}

// Generate creates a signed token for the given user that starts at the given
// point in time.
func (auth *tokenAuth) Generate(user *User, now time.Time) (string, error) {
	if user == nil || user.Name == "" {
		return "", errInvalidToken
	}
	claims, err := json.Marshal(tokenClaims{
		Subject: user.Name,
		Issued:  now.Unix(),
		Expires: now.Add(tokenTTL).Unix(),
	})
	if err != nil {
		return "", err
	}
	body := base64url([]byte(`{"alg":"HS256","typ":"JWT"}`)) + "." +
		base64url(claims)
	return body + "." + auth.sign(body), nil
}

// Verify checks the signature and the expiration of the given token and
// returns the user name that is stored in it.
func (auth *tokenAuth) Verify(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errInvalidToken
	}
	body := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(auth.sign(body))) {
		return "", errInvalidToken
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errInvalidToken
	}
	claims := &tokenClaims{}
	if err := json.Unmarshal(data, claims); err != nil {
		return "", errInvalidToken
	}
	if claims.Subject == "" {
		return "", errInvalidToken
	}
	if claims.Expires > 0 && time.Now().Unix() > claims.Expires {
		return "", errExpiredToken
	}
	return claims.Subject, nil
}

// sign returns the signature of the given token body.
func (auth *tokenAuth) sign(body string) string {
	mac := hmac.New(sha256.New, auth.key)
	mac.Write([]byte(body))
	return base64url(mac.Sum(nil))
}

func base64url(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
