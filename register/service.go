package register

import (
	"context"
	"github.com/One-Regular-Guy/shop-page-api/database"
	"github.com/redis/go-redis/v9"
	"log"
)

type Service interface {
	Cache() *redis.Client
	DB() *database.Queries
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	RegisterUser(name, username, email, password string) error
}

type service struct {
	cache *redis.Client
	db    *database.Queries
}

func (s *service) Cache() *redis.Client {
	return s.cache
}
func (s *service) DB() *database.Queries {
	return s.db
}

func NewService(cache *redis.Client, db *database.Queries) Service {
	return &service{
		cache: cache,
		db:    db,
	}
}

func (s *service) CheckUsername(username string) bool {
	// Primeiro verifica no Redis
	val, err := s.cache.Get(context.Background(), "username:"+username).Result()
	if err == nil && val != "" {
		log.Printf("Username found in Redis")
		return true
	}
	// Se não encontrar no Redis, verifica no banco de dados
	value, err := s.db.GetUsername(context.Background(), username)
	if err != nil {
		log.Printf("Error getting username: %v", err)
		return false
	}
	if value == "" {
		log.Printf("Username not found")
		return false
	}
	log.Printf("Username found in DB")
	return true
}

func (s *service) CheckEmail(email string) bool {
	// Primeiro verifica no Redis
	val, err := s.cache.Get(context.Background(), "email:"+email).Result()
	if err == nil && val == "exists" {
		log.Printf("Email found in Redis")
		return true
	}
	// Se não encontrar no Redis, verifica no banco de dados
	_, err = s.db.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("Error getting email: %v", err)
		return false
	}
	status := s.Cache().Set(context.Background(), "email:"+email, "exists", 0)
	if status.Err() != nil {
		log.Printf("Error setting email status in Redis: %v", status.Err())
		return true
	}
	log.Printf("Email found in DB")
	return true
}

func (s *service) RegisterUser(name, username, email, password string) error {
	err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	})
	return err
}

var ServiceInstance Service
