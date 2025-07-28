package infrastructure

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"os"
	"time"
)

type JWTTokenProvider struct {
	privKey jwk.Key
	pubKey  jwk.Key
}

func NewJWTTokenProvider() (*JWTTokenProvider, error) {
	privFile, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privFile)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, err
	}
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	jwkprivkey, err := jwk.Import(privkey)
	if err != nil {
		return nil, err
	}
	pubkey, err := jwk.PublicKeyOf(jwkprivkey)
	if err != nil {
		return nil, err
	}
	return &JWTTokenProvider{privKey: jwkprivkey, pubKey: pubkey}, nil
}

func (p *JWTTokenProvider) GenerateToken(username string) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer("shop-page-api").
		Subject(username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), p.privKey))
	if err != nil {
		return "", err
	}
	encrypted, err := jwe.Encrypt(signed, jwe.WithKey(jwa.RSA_OAEP(), p.pubKey))
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
