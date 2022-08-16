package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Get
	r.GET("/accounts/:id", func(c *gin.Context) {
		acct_id := c.Params.ByName("id")
		value, ok := db[acct_id] // call CRDB
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Security 
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"user0": "pw0" // user & password
		"user1": "pw1" // user & password
	}))

	// Post
	r.POST("accounts", func(c *gin.Context) {
		var json struct {
			Value string `json:"value" binding:"required"`
		}
		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}