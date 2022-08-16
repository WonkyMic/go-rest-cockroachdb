package main

import (
	"context"
	// "fmt"
	"log"
	// "os"
	// "time"

	"github.com/jackc/pgx/v4"
	"github.com/google/uuid"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
)

func insert(ctx context.Context, tx pgx.Tx, accts [4]uuid.UUID) error {
	log.Println("main.insert")
	if _, err := tx.Exec(ctx,
        "INSERT INTO accounts (id, balance) VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8)",
		 accts[0], 250, accts[1], 100, accts[2], 500, accts[3], 300); err != nil {
		return err
	}
	return nil
}

func printBalances(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "SELECT id, balance FROM accounts")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id uuid.UUID
        var balance int
        if err := rows.Scan(&id, &balance); err != nil {
            log.Fatal(err)
        }
        log.Printf("%s: %d\n", id, balance)
    }
    return nil
}

func transferFunds(ctx context.Context, tx pgx.Tx, from uuid.UUID, to uuid.UUID, amount int) error {
    // Read the balance.
    var fromBalance int
    if err := tx.QueryRow(ctx,
        "SELECT balance FROM accounts WHERE id = $1", from).Scan(&fromBalance); err != nil {
        return err
    }

    if fromBalance < amount {
        log.Println("insufficient funds")
    }

    // Perform the transfer.
    log.Printf("Transferring funds from account with ID %s to account with ID %s...", from, to)
    if _, err := tx.Exec(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from); err != nil {
        return err
    }
    if _, err := tx.Exec(ctx,
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to); err != nil {
        return err
    }
    return nil
}

func deleteRows(ctx context.Context, tx pgx.Tx, one uuid.UUID, two uuid.UUID) error {
    // Delete two rows into the "accounts" table.
    log.Printf("Deleting rows with IDs %s and %s...", one, two)
    if _, err := tx.Exec(ctx,
        "DELETE FROM accounts WHERE id IN ($1, $2)", one, two); err != nil {
        return err
    }
    return nil
}

func main() {

	// Create connection
	dsn := "postgresql://wonkymic:<password>@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	defer conn.Close(context.Background()) // defer :: will execute when surrounding function returns
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// Insert initial rows
    var accounts [4]uuid.UUID
    for i := 0; i < len(accounts); i++ {
        accounts[i] = uuid.New()
    }

    //TODO :: implement this
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
        return insert(context.Background(), tx, accounts)
    })
    if err == nil {
        log.Println("New rows created.")
    } else {
        log.Fatal("error: ", err)
    }

	log.Println("Init:")
	printBalances(conn)
}