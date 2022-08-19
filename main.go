package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"wonkymic/go-service-crdb/domain"
	"wonkymic/go-service-crdb/data"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	// Ping Pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// GET :: User
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		user := data.GetUser(dbpool, id)
		c.JSON(http.StatusOK, user)
	})
	// Delete :: User
	r.DELETE("/user/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		data.DeleteUser(dbpool, id)
		c.JSON(http.StatusOK, gin.H{"message": "User Deleted"})
	})
	// GET :: All Users
	r.GET("/user", func(c *gin.Context) {
		users := data.GetUsers(dbpool)
		c.JSON(http.StatusOK, users)
	})	
	// CREATE :: User
	r.POST("/user", func(c *gin.Context) {
		var user_req domain.UserReq
		c.BindJSON(&user_req)
		user_res := data.CreateUser(dbpool, user_req)
		c.JSON(http.StatusCreated, user_res)
	})
	return r
}

func main() {
	// Set connection pool configuration 
	database_url := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(database_url)
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	// Create connection pool to CRDB
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer dbpool.Close() // close connection pool when main exits

	r := setupRouter(dbpool)
	r.Run(":8080")
}