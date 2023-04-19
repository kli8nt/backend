package controllers

import (
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
	"github.com/gin-gonic/gin"
)

func GetRepo(c *gin.Context) {
	var repo models.Repo
	initializers.DB.First(&repo, c.Param("id"))
	c.JSON(200, gin.H{"repo": repo})
}

func FetchReposByUser(c *gin.Context) {
	var repos []models.Repo
	initializers.DB.Find(&repos, "Owner = ?", c.Param("id"))
	c.JSON(200, gin.H{"repos": repos})
}

func AddUserRepository(repo models.Repo) {
	initializers.DB.Create(&repo)
}
