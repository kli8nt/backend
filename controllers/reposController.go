package controllers

import (
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
	"github.com/gin-gonic/gin"
)

func GetRepo(c *gin.Context) {
	var repo models.Repo
	var owner models.User
	initializers.DB.Where("username = ?", c.Param("username")).First(&owner)
	initializers.DB.First(&repo, "owner_id = ? AND name = ?", owner.ID, c.Param("name"))
	c.JSON(200, gin.H{"repo": repo})
}

func FetchReposByUser(c *gin.Context) {
	var repos []models.Repo
	var owner models.User
	initializers.DB.Where("username = ?", c.Param("username")).First(&owner)
	initializers.DB.Find(&repos, "owner_id = ?", owner.ID)
	c.JSON(200, gin.H{"repos": repos})
}

func AddUserRepository(repo models.Repo) {
	initializers.DB.Create(&repo)
}
