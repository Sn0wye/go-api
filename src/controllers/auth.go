package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Sn0wye/go-api/pkg/exceptions"
	"github.com/Sn0wye/go-api/pkg/jwt"
	"github.com/Sn0wye/go-api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController interface {
	Profile(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GenerateToken(ctx context.Context, userId string) (string, error)
}

type authController struct {
	db  *gorm.DB
	jwt *jwt.JWT
}

func NewAuthController(db *gorm.DB, jwt *jwt.JWT) AuthController {
	return &authController{db: db, jwt: jwt}
}

func (s *authController) Profile(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.MyCustomClaims)
	user := models.User{}
	s.db.Where("id = ?", claims.UserId).First(&user)

	c.JSON(http.StatusOK, user)
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *authController) Register(c *gin.Context) {
	db := s.db
	body := RegisterRequest{}
	err := c.ShouldBindJSON(&body)
	var user models.User

	if err != nil {
		exceptions.UnprocessableEntity(c, "Invalid JSON provided")
		return
	}

	exists := db.Where("email = ?", body.Email).First(&user).RowsAffected
	if exists > 0 {
		exceptions.BadRequest(c, "Email already taken")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		exceptions.InternalServerError(c, "failed to hash password")
		return
	}

	user = models.User{
		Username: body.Username,
		Password: string(hashedPassword),
		Email:    body.Email,
		Name:     body.Name,
	}

	db.Create(&user)

	token, err := s.GenerateToken(c, user.ID.String())
	if err != nil {
		exceptions.InternalServerError(c, "failed to generate JWT token")
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *authController) Login(c *gin.Context) {
	db := s.db
	body := LoginRequest{}
	err := c.ShouldBindJSON(&body)
	var user models.User

	if err != nil {
		exceptions.UnprocessableEntity(c, "Invalid JSON provided")
		return
	}

	db.Where("email = ?", body.Email).First(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		exceptions.Unauthorized(c)
		return
	}

	token, err := s.GenerateToken(c, user.ID.String())
	if err != nil {
		exceptions.InternalServerError(c, "failed to generate JWT token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *authController) GenerateToken(ctx context.Context, userId string) (string, error) {
	token, err := s.jwt.GenToken(userId, time.Now().Add(time.Hour*24*90)) // 90 days
	if err != nil {
		return "", errors.Wrap(err, "failed to generate JWT token")
	}

	return token, nil
}
