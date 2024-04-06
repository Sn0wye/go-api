package middleware

import (
	"fmt"
	"time"

	"github.com/Sn0wye/go-api/pkg/exceptions"
	"github.com/Sn0wye/go-api/pkg/jwt"
	"github.com/Sn0wye/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func JWTMiddleware(conf *viper.Viper, logger *logger.Logger) gin.HandlerFunc {
	j := jwt.NewJwt(conf)
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			fmt.Println("No token provided")
			exceptions.Unauthorized(ctx)
			ctx.Abort()
			return
		}

		fmt.Print(tokenString)

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			fmt.Println("Invalid token provided")
			exceptions.Unauthorized(ctx)
			ctx.Abort()
			return
		}

		expirationTime := claims.ExpiresAt.Time
		if time.Until(expirationTime) < 5*time.Minute {
			newTokenString, err := j.GenToken(claims.UserId, time.Now().Add(time.Hour*24*90))
			if err != nil {
				fmt.Println("Error generating new token")
				exceptions.InternalServerError(ctx, err.Error())
				ctx.Abort()
				return
			}
			ctx.Header("Authorization", "Bearer "+newTokenString)
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *logger.Logger) {
	userInfo := ctx.MustGet("claims").(*jwt.MyCustomClaims)
	logger.NewContext(ctx, zap.String("UserId", userInfo.UserId))
}
