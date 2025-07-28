package infrastructure

import (
	"context"
	"github.com/One-Regular-Guy/shop-page-api/internal/domain"
	database2 "github.com/One-Regular-Guy/shop-page-api/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	cache *redis.Client
	db    *database2.Queries
}

func NewUserRepository(cache *redis.Client, db *database2.Queries) *UserRepositoryImpl {
	return &UserRepositoryImpl{cache: cache, db: db}
}

func (r *UserRepositoryImpl) IsUsernameTaken(username string) (bool, error) {
	val, err := r.cache.Get(context.Background(), "username:"+username).Result()
	if err == nil && val != "" {
		return true, nil
	}
	value, err := r.db.GetUsername(context.Background(), username)
	if err != nil || value == "" {
		return false, nil
	}
	return true, nil
}

func (r *UserRepositoryImpl) IsEmailTaken(email string) (bool, error) {
	val, err := r.cache.Get(context.Background(), "email:"+email).Result()
	if err == nil && val == "exists" {
		return true, nil
	}
	_, err = r.db.GetUserByEmail(context.Background(), email)
	if err != nil {
		return false, nil
	}
	_ = r.cache.Set(context.Background(), "email:"+email, "exists", 0)
	return true, nil
}

func (r *UserRepositoryImpl) CreateUser(user *domain.User) error {
	id := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.db.CreateUser(context.Background(), database2.CreateUserParams{
		ID:       uuid.MustParse(id),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
	})
}

func (r *UserRepositoryImpl) CheckPassword(username, password string) (bool, error) {
	storedPassword, err := r.db.GetPassword(context.Background(), username)
	if err != nil || storedPassword == "" {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

var _ domain.UserRepository = (*UserRepositoryImpl)(nil)
