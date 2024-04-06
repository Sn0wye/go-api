package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(c *gin.Context, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	c.JSON(code, resp)
}

var (
	BadRequest = func(c *gin.Context, message string) {
		writeError(c, message, http.StatusBadRequest)
	}
	InternalServerError = func(c *gin.Context, message string) {
		writeError(c, message, http.StatusInternalServerError)
	}
	Unauthorized = func(c *gin.Context) {
		writeError(c, "Unauthorized", http.StatusUnauthorized)
	}
	UnprocessableEntity = func(c *gin.Context, message string) {
		writeError(c, message, http.StatusUnprocessableEntity)
	}
)
