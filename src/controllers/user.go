package controllers

import (
	"net/http"

	"github.com/Sn0wye/go-api/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &userController{db: db}
}

type CreateUser struct {
	Name string `json:"name"`
}

func (uc *userController) CreateUser(c *gin.Context) {
	var user CreateUser
	c.BindJSON(&user)

	newUser := models.User{
		Name: user.Name,
	}

	uc.db.Create(&newUser)

	c.JSON(http.StatusCreated, newUser)
}

func (uc *userController) GetUsers(c *gin.Context) {
	var users []models.User

	uc.db.Find(&users)

	c.JSON(http.StatusOK, users)
}

func (uc *userController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := uc.db.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	uc.db.Delete(&user)

	c.JSON(http.StatusOK, user)
}
