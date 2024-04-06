package main

import (
	"fmt"
	"net/http"

	"github.com/Sn0wye/go-api/pkg/config"
	"github.com/Sn0wye/go-api/pkg/logger"
	"github.com/Sn0wye/go-api/src/db"
	"github.com/Sn0wye/go-api/src/middleware"
	"github.com/Sn0wye/go-api/src/migration"
	"github.com/Sn0wye/go-api/src/models"
	"github.com/Sn0wye/go-api/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.GetConfig()
	logger := logger.NewLog(conf)
	db := db.GetDB()

	db.AutoMigrate(
		models.RetrieveAll()...,
	)

	migrate := migration.NewMigrate(
		db,
		logger,
	)

	migrate.Run()

	jwt := middleware.JWTMiddleware(conf, logger)

	router := gin.Default()
	v1 := router.Group("/v1")
	v1.Use(middleware.CORSMiddleware())

	v1.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Version 1.0",
		})
	})

	routes.BindAuthRoutes(v1, jwt, logger)
	routes.BindUserRoutes(v1, jwt, logger)

	port := conf.Get("http.port")
	formattedPort := fmt.Sprintf(":%d", port)
	fmt.Printf("Server is running on port %d\n", port)
	http.ListenAndServe(formattedPort, router)
}
