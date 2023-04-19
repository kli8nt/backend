package controllers

import (
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
	"github.com/gin-gonic/gin"
)

func AddUser(user models.User) {
	initializers.DB.Create(&user)
}

func GetUser(c *gin.Context) {
	var user models.User
	initializers.DB.First(&user, "username = ?", c.Param("username"))
	c.JSON(200, gin.H{"user": user})
}

func UserExists(username string) bool {
	var user models.User
	initializers.DB.First(&user, "username = ?", username)
	return user.Username != ""
}

func UpdateUserToken(username string, token string) {
	initializers.DB.Model(&models.User{}).Where("username = ?", username).Update("token", token)
}
