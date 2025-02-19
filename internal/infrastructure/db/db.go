package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func ConnectDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		panic(err)
	}

	slog.Info("Database connected successfully!")
	return db
}
