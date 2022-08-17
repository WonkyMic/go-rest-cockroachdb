package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"wonkymic/go-service-crdb/domain"
	// "wonkymic/go-service-crdb/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// User
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var user_res domain.UserRes
		err := dbpool.QueryRow(context.Background(),
			"SELECT id, name FROM users WHERE id = $1", id).Scan(&user_res.Id, &user_res.Name);

		if err == nil {
			c.JSON(http.StatusOK, user_res)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		}
	})
	// GET :: All Users
	r.GET("/user", func(c *gin.Context) {
		rows, err := dbpool.Query(context.Background(), "select * from public.users")
		
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
	// GET :: All Accounts
	r.GET("/account", func(c *gin.Context) {
		rows, err := dbpool.Query(context.Background(), "select * from public.accounts")
		
		if err != nil {
			log.Println("error while executing select accounts query")
		}

		var id uuid.UUID
		var balance int
		for rows.Next() {
			if err := rows.Scan(&id, &balance); err!= nil {
				log.Fatal("error while iterating dataset")
			}
			log.Println("[id:", id, "balance:", balance, "]")
		}
	})

	r.POST("/user", func(c *gin.Context) {
		q := "INSERT INTO users(id, name) VALUES($1, $2) RETURNING id, name"
		// level := pgx.Serializable

		// Map request body
		var user_req domain.UserReq
		c.BindJSON(&user_req)

		// // CRDB	transaction
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

		c.JSON(200, user_res)
	})
	return r
}

func main() {
	// Create connection
	database_url := os.Getenv("DATABASE_URL")

	// Set connection pool configuration, with max connection pool size.
	config, err := pgxpool.ParseConfig(database_url)
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	// Create a connecction pool to the "bank" database
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer dbpool.Close()

	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	r := setupRouter(dbpool)
	r.Run(":8080")
}