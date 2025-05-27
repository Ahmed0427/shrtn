package server

import (
	"net/http"
	"database/sql"
	"log"

	"github.com/ahmed0427/shrtn/internal/db"
	"github.com/ahmed0427/shrtn/internal/utils"

	_ "github.com/lib/pq"
)

const CACHE_CAPACITY = 1000

type Config struct {
	db  *db.Queries
	cache *utils.LRUCache
	cacheHits int
}

func NewRouter(connStr string) *http.ServeMux {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := db.New(dbConn)
	cfg := Config {
		db: dbQueries,
		cache: utils.NewLRUCache(CACHE_CAPACITY),
		cacheHits: 0,
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /", cfg.handleShortening)
	router.HandleFunc("GET /{id}", cfg.handleRedirection)
	return router
}
