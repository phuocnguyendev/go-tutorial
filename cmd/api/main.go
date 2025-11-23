package main

import (
	"fmt"
	httpdelivery "go-tutorial/internal/delivery/http"
	"log"
	"os"

	"go-tutorial/internal/domain"
	"go-tutorial/internal/infrastructure/database"
	"go-tutorial/internal/repository/postgresrepo"
	"go-tutorial/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Swagger docs
	// import generated docs to register swagger spec
	_ "go-tutorial/internal/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Tutorial API
// @version 1.0
// @description API Documentation for Go Tutorial project
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// 0. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 1. Kết nối DB
	db := database.NewPostgresDB()

	// 2. Auto-migrate
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("auto migrate failed:", err)
	}

	// 3. Khởi tạo repository + usecase
	userRepo := postgresrepo.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo)

	// 4. Khởi tạo Gin
	r := gin.Default()

	// 5. Đăng ký route API
	httpdelivery.NewAuthHandler(r, authUC)

	// ⭐ Đăng ký Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Debug env
	fmt.Println("POSTGRES:", os.Getenv("POSTGRES_USER"))
	fmt.Println("JWT:", os.Getenv("JWT_SECRET"))

	// 6. Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
