package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GithubHooks(c *gin.Context) (string, string) {

	// extract username, reponame from the request
	type Payload struct {
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	var p Payload

	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Panic(err)
	}
	username := p.Repository.Owner.Login
	reponame := p.Repository.Name

	return username, reponame

}
