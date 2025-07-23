package main

import (
	"context"
	"github.com/One-Regular-Guy/shop-page-api/database"
	"github.com/One-Regular-Guy/shop-page-api/login"
	"github.com/One-Regular-Guy/shop-page-api/register"
	"github.com/One-Regular-Guy/shop-page-api/status"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	databaseUrl := os.Getenv("DATABASE_URL")
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if databaseUrl == "" || redisUrl == "" || redisPassword == "" {
		log.Fatal("Environment variables DATABASE_URL, REDIS_URL, and REDIS_PASSWORD must be set")
	}
	if err != nil {
		log.Print("Error loading .env file")
	}
	databaseCtx := context.Background()
	pool, err := pgxpool.New(databaseCtx, databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	db := database.New(pool)
	cache := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	defer func(cache *redis.Client) {
		err := cache.Close()
		if err != nil {
			log.Print(err)
		}
	}(cache)
	login.ServiceInstance = login.NewService(cache, db)
	register.ServiceInstance = register.NewService(cache, db)
	log.Print("Defining Server Type")
	router := gin.Default()
	log.Print("Registering endpoints ...")
	router.GET("/hello", status.Handler)
	router.POST("/login", login.Handler)
	router.POST("/register", register.Handler)
	log.Print("Endpoints registred ...")
	log.Print("Starting server on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Server finished Gracefully")
}
