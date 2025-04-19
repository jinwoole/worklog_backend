// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jinwoole/worklog-backend/config"
	"github.com/jinwoole/worklog-backend/handler"
	"github.com/jinwoole/worklog-backend/middleware"
	"github.com/jinwoole/worklog-backend/repository"
	"github.com/jinwoole/worklog-backend/service"
)

func main() {
	// Load .env for PORT and JWT_SECRET
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	// Initialize Database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	// main.go (InitDB 직후)
	rows, err := db.Query("SELECT schemaname, tablename FROM pg_tables WHERE schemaname NOT IN ('pg_catalog','information_schema')")
	if err != nil {
		fmt.Print("failed to list tables: %v", err)
	}
	fmt.Println("🚀 Database connection established successfully!")
	defer rows.Close()
	for rows.Next() {
		var schema, table string
		rows.Scan(&schema, &table)
		fmt.Printf("found table → %s.%s\n", schema, table)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	workLogRepo := repository.NewWorkLogRepository(db)

	// Initialize services
	userSvc := service.NewUserService(userRepo)
	workLogSvc := service.NewWorkLogService(workLogRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userSvc)
	workLogHandler := handler.NewWorkLogHandler(workLogSvc)

	// Setup Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	router.POST("/api/register", userHandler.Register)
	router.POST("/api/login", userHandler.Login)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "🚀 WorkLog API server is running!")
	})

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/worklog", workLogHandler.CreateWorkLog)
		api.PUT("/worklog", workLogHandler.UpdateWorkLog)
		api.GET("/worklog", workLogHandler.GetAllWorkLogs)
		api.GET("/me", workLogHandler.GetMe)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
