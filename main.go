package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"wonkymic/go-service-crdb/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	// Ping Pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// GET :: User
	r.GET("/user/:id", func(c *gin.Context) {
		q := "SELECT id, name FROM users WHERE id = $1"
		id := c.Params.ByName("id")
		var user_res domain.UserRes
		err := dbpool.QueryRow(context.Background(), q, id).Scan(&user_res.Id, &user_res.Name);

		if err == nil {
			c.JSON(http.StatusOK, user_res)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		}
	})
	// Delete :: User
	r.DELETE("/user/:id", func(c *gin.Context) {
		q := "DELETE FROM users WHERE id = $1"
		id := c.Params.ByName("id")
		// CRDB	transaction
		ctx := context.Background()
		// TODO - update Isolation level (pgx.TxOptions{})
		tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			log.Fatal("error: ", err)
		}

		// Map Response
		_, err = tx.Exec(ctx, q, id)
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("error: ", err)
		}

		c.JSON(200, gin.H{"message": "User Deleted"})
	})
	// GET :: All Users
	r.GET("/user", func(c *gin.Context) {
		q := "SELECT * FROM users"
		rows, err := dbpool.Query(context.Background(), q)
		
		if err != nil {
			log.Println("error while executing select users query")
		}

		users := make([]domain.UserRes, 0)
		for rows.Next() {
			var user_res domain.UserRes
			if err := rows.Scan(&user_res.Id, &user_res.Name); err!= nil {
				log.Fatal("error while iterating dataset")
			}
			users = append(users, user_res)
		}
		if len(users) > 0 {
			c.JSON(http.StatusOK, users)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Users Not Found"})
		}
	})	
	// CREATE :: User
	r.POST("/user", func(c *gin.Context) {
		q := "INSERT INTO users(id, name) VALUES($1, $2) RETURNING id, name"

		// Map request body
		var user_req domain.UserReq
		c.BindJSON(&user_req)

		// CRDB	transaction
		ctx := context.Background()
		// TODO - update Isolation level (pgx.TxOptions{})
		tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			log.Fatal("error: ", err)
		}

		// Map Response
		var user_res domain.UserRes
		err = tx.QueryRow(ctx, q, uuid.New(), user_req.Name).Scan(&user_res.Id, &user_res.Name)
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("error: ", err)
		}

		c.JSON(201, user_res)
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