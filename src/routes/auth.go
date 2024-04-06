package routes

import (
	"github.com/Sn0wye/go-api/pkg/config"
	"github.com/Sn0wye/go-api/pkg/jwt"
	"github.com/Sn0wye/go-api/pkg/logger"
	"github.com/Sn0wye/go-api/src/controllers"
	"github.com/Sn0wye/go-api/src/db"
	"github.com/gin-gonic/gin"
)

func BindAuthRoutes(r *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, log *logger.Logger) {
	db := db.GetDB()
	conf := config.GetConfig()
	jwt := jwt.NewJwt(conf)

	router := r.Group("/auth")
	controller := controllers.NewAuthController(db, jwt)

	router.POST("/login", controller.Login)
	router.POST("/register", controller.Register)

	router.Use(jwtMiddleware)
	router.GET("/profile", controller.Profile)
}
