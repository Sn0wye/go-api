package routes

import (
	"github.com/Sn0wye/go-api/pkg/logger"
	"github.com/Sn0wye/go-api/src/controllers"
	"github.com/Sn0wye/go-api/src/db"
	"github.com/gin-gonic/gin"
)

func BindUserRoutes(r *gin.RouterGroup, jwt gin.HandlerFunc, logger *logger.Logger) {
	db := db.GetDB()
	controller := controllers.NewUserController(db)
	router := r.Group("/users")
	router.Use(jwt)

	router.GET("", controller.GetUsers)
	router.POST("", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
}
