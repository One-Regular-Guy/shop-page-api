package main

import (
	"context"
	"github.com/One-Regular-Guy/shop-page-api/internal/handler"
	"github.com/One-Regular-Guy/shop-page-api/internal/infrastructure"
	"github.com/One-Regular-Guy/shop-page-api/internal/infrastructure/database"
	"github.com/One-Regular-Guy/shop-page-api/internal/usecase"
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
		panic("Unable to connect to database: " + err.Error())
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
	// Remover instâncias antigas de login e register
	// login.ServiceInstance = login.NewService(cache, db)
	// register.ServiceInstance = register.NewService(cache, db)
	// --- Clean Architecture: Registro ---
	userRepo := infrastructure.NewUserRepository(cache, db)
	registerUseCase := &usecase.RegisterUserUseCase{UserRepo: userRepo}
	registerHandler := &handler.RegisterHandler{UseCase: registerUseCase}
	// ---
	// --- Clean Architecture: Login ---
	jwtProvider, err := infrastructure.NewJWTTokenProvider()
	if err != nil {
		log.Fatalf("Failed to initialize JWT provider: %v", err)
	}
	loginUseCase := &usecase.LoginUserUseCase{UserRepo: userRepo, TokenProvider: jwtProvider}
	loginHandler := &handler.LoginHandler{UseCase: loginUseCase}
	// ---
	router := gin.Default()
	// Middleware CORS simples
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.POST("/login", loginHandler.Handle)
	router.OPTIONS("/login", func(c *gin.Context) { c.Status(204) })
	router.POST("/register", registerHandler.Handle)
	router.OPTIONS("/register", func(c *gin.Context) { c.Status(204) })
	err = router.Run(":8084")
	if err != nil {
		log.Fatal(err)
	}
}
