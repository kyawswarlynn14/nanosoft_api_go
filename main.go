package main

import (
	"log"
	"nanosoft/middleware"
	"nanosoft/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	publicRoutes := router.Group("/")

	authenticatedRoutes := router.Group("/")
	authenticatedRoutes.Use(middleware.Authentication())

	adminRoutes := router.Group("/")
	adminRoutes.Use(middleware.Authentication())
	adminRoutes.Use(middleware.AuthorizeRole([]int{1, 2}))

	routes.UserRoutes(publicRoutes, authenticatedRoutes, adminRoutes)
	routes.ServiceRoutes(publicRoutes, authenticatedRoutes, adminRoutes)
	routes.ProjectRoutes(publicRoutes, authenticatedRoutes, adminRoutes)
	routes.RemarkRoutes(publicRoutes, authenticatedRoutes, adminRoutes)
	routes.EmailRoutes(publicRoutes, authenticatedRoutes, adminRoutes)
	log.Fatal(router.Run(":" + port))
}
