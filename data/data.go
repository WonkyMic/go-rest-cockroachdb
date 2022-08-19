package data

import (
	"context"
	"log"
	"wonkymic/go-service-crdb/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func CreateUser(dbpool *pgxpool.Pool, user domain.UserReq) domain.UserRes {
	q := "INSERT INTO users(id, name) VALUES($1, $2) RETURNING id, name"
	ctx := context.Background()
	// TODO - update Isolation level (pgx.TxOptions{})
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})

	// Map Response
	var user_res domain.UserRes
	err = tx.QueryRow(ctx, q, uuid.New(), user.Name).Scan(&user_res.Id, &user_res.Name)
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return user_res
}

func GetUsers(dbpool *pgxpool.Pool) []domain.UserRes {
	q := "SELECT id, name FROM users"
	rows, err := dbpool.Query(context.Background(), q)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer rows.Close()
	var users []domain.UserRes
	for rows.Next() {
		var user domain.UserRes
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Fatal("error: ", err)
		}
		users = append(users, user)
	}
	return users
}

func GetUser(dbpool *pgxpool.Pool, id string) domain.UserRes {
	q := "SELECT id, name FROM users WHERE id = $1"
	var user domain.UserRes
	err := dbpool.QueryRow(context.Background(), q, id).Scan(&user.Id, &user.Name)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return user
}

func DeleteUser(dbpool *pgxpool.Pool, id string) {
	q := "DELETE FROM users WHERE id = $1"
	ctx := context.Background()
	// TODO - update Isolation level (pgx.TxOptions{})
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
}