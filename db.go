package main

import (
	"context"
	"log"
	"refeclt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kjrodgers26/solid-api/db"
)

func GerUsers() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "") // TODO: added postgre connection string
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	users, err := queries.ListUsers()
	if err != nil {
		return err
	}
}
