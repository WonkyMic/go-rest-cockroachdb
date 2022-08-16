package data

import (
	"wonkymic/go-service-crdb/domain"

	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v4"
	// "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
)

func Insert_user(ctx context.Context, tx pgx.Tx, user domain.UserReq) error {
	log.Println("crdb.Insert")
	fmt.Printf("%+V", user)
	// if _, err := tx.Exec(ctx,
	// 	"INSERT INTO users(id, name) VALUES($1, $2)"),
	// 	user.Id, user.Name); err != nil {
	// 	return nil
	// }
	return nil

}

func Insert_user_without_print(ctx context.Context, tx pgx.Tx, user domain.UserReq) error {
	log.Println("crdb.Insert")
	fmt.Printf("%+V", user)
	if _, err := tx.Exec(ctx,
		"INSERT INTO users(id, name) VALUES($1, $2)",
		user.Id, user.Name); err != nil {
		return err
	}
	return nil
}