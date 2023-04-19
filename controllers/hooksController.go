package controllers

import "github.com/gin-gonic/gin"

func GithubHooks(c *gin.Context) {
	c.JSON(200, gin.H{"message": "got a hook!"})
}
