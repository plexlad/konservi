package main

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/plexlad/konservi/ent"
)

var db *ent.Client

func InitDB(dsn string) error {
	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	db = client
	return nil
}
