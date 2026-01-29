package main

import (
	"app/config"
	"app/internal/api"
	"app/internal/database"
	"app/internal/repository"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Template struct {
	templates *template.Template
}

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	db := database.ConnectDatabase()

	userRepository := repository.NewUserRepository(db)
	userHandler := api.NewUserHandler(userRepository)

	router := gin.Default()
	router.GET("/health", CheckHealth)
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUserById)
	router.GET("/users", userHandler.GetUsers)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // default fallback
	}
	slog.Info("server started successfully", "port", port)
	if err := router.Run(port); err != nil {
		slog.Error("failed to start server", "error", err)
		panic(err)
	}
}

func CheckHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Service running successfully"})
}
