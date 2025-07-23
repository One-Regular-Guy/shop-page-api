package login

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/One-Regular-Guy/shop-page-api/database"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

// Defines methods that the Service should implement
type Service interface {
	Cache() *redis.Client
	DB() *database.Queries
	PrivKey() jwk.Key
	PubKey() jwk.Key
	CheckUsername(username string) bool
	CheckPassword(username, incomingPassword string) bool
	ProvideToken(username string) (string, error)
}

// Atributes of the Service
type service struct {
	cache   *redis.Client
	db      *database.Queries
	privKey jwk.Key
	pubKey  jwk.Key
}

// Getters for the service struct
func (s *service) Cache() *redis.Client {
	return s.cache
}
func (s *service) DB() *database.Queries {
	return s.db
}
func (s *service) PrivKey() jwk.Key {
	return s.privKey
}
func (s *service) PubKey() jwk.Key {
	return s.pubKey
}

// Constructor for the Service
func NewService(cache *redis.Client, db *database.Queries) Service {
	privFile, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Printf("failed to read private key file: %s\n", err)
		panic(err)
	}
	block, _ := pem.Decode(privFile)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("Failed to decode PEM block containing private key")
	}

	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("failed to parse private key: %s\n", err)
		panic(err)
	}

	jwkprivkey, err := jwk.Import(privkey)
	if err != nil {
		log.Printf("failed to parse JWK: %s\n", err)
		panic(err)
	}

	pubkey, err := jwk.PublicKeyOf(jwkprivkey)
	if err != nil {
		log.Printf("failed to get public key: %s\n", err)
		panic(err)
	}
	return &service{
		cache:   cache,
		db:      db,
		privKey: jwkprivkey,
		pubKey:  pubkey,
	}
}

// Methods for the Service

func (s *service) CheckUsername(username string) bool {
	value, err := s.db.GetUsername(context.Background(), username)
	if err != nil {
		log.Printf("Error getting username: %v", err)
		return false
	}
	if value == "" {
		log.Printf("Username not found")
		return false
	}
	log.Printf("Username found")
	return true
}
func (s *service) CheckPassword(username, incomingPassword string) bool {
	storedPassword, err := s.DB().GetPassword(context.Background(), username)
	if err != nil {
		log.Print("Error getting password for user")
		return false
	}
	if storedPassword == "" {
		log.Print("No password found for user")
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(incomingPassword))
	if err != nil {
		log.Print("Password mismatch for user")
		return false
	}
	log.Print("Password match for user")
	return true
}
func (s *service) ProvideToken(username string) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer("shop-page-api").
		Subject(username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), s.PrivKey()))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return "", err
	}
	encrypted, err := jwe.Encrypt(signed, jwe.WithKey(jwa.RSA_OAEP(), s.PubKey()))
	if err != nil {
		fmt.Printf("failed to encrypt payload: %s\n", err)
		return "", err
	}
	return string(encrypted), nil
}

// Holds necessary dependencies for the service
var ServiceInstance Service
