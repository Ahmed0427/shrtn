package server

import (
	"net/http"
	"database/sql"
	"log"

	"github.com/ahmed0427/shrtn/internal/db"

	_ "github.com/lib/pq"
)

type Config struct {
	db  *db.Queries
}

func NewRouter(connStr string) *http.ServeMux {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := db.New(dbConn)
	cfg := Config {
		db: dbQueries,
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /", cfg.handleShortening)
	router.HandleFunc("GET /{id}", cfg.handleRedirection)
	return router
}
