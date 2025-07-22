package login

import (
	"context"
	"github.com/One-Regular-Guy/shop-page-api/database"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Cache() *redis.Client
	DB() *database.Queries
	CheckUsername(username string) error
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

func (s *service) CheckUsername(username string) error {
	_, err := s.db.GetPassword(context.Background(), username)
	if err != nil {
		return err
	}
	return nil
}

var ServiceInstance Service
