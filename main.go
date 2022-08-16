package main

import (
	"context"
	"net/http"
	"wonkymic/go-service-crdb/domain"
	"wonkymic/go-service-crdb/data"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
)

var db = make(map[string]string)

func setupRouter(ctx context.Context, tx pgx.Tx) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// User
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user] // call CRDB
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})
	r.POST("/user", func(c *gin.Context) {
		var user_req domain.UserReq

		c.BindJSON(&user_req)

		// CRDB	transaction
		err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return data.Insert_user(ctx, tx, *user_req)
		})
		if err == nil {
			log.Println("New user created.")
		} else {
			log.Fatal("error: ", err)
		}
		c.JSON(200, user_req)
	})
	return r
}

func main() {
	// Create connection
	// TODO :: connection pooling
	
	// Set connection pool configuration, with max connection pool size.
	config, err := pgxpool.ParseConfig("postgresql://wonkymic:AEOLKPUjsaDjYS2UYDe60w@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43")
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	// Create a connecction pool to the "bank" database
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer dbpool.Close()

	dsn := "postgresql://wonkymic:<password>@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	defer conn.Close(context.Background()) // defer :: will execute when surrounding function returns
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	r := setupRouter(context.Background(), conn)
	r.Run(":8080")
}