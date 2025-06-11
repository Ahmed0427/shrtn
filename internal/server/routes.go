package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ahmed0427/shrtn/internal/db"
	"github.com/ahmed0427/shrtn/internal/utils"

	_ "github.com/lib/pq"
)

const CACHE_CAPACITY = 1000

type Config struct {
	db        *db.Queries
	cache     *utils.LRUCache
	cacheHits int
	visitors  visitors
}

func NewRouter(connStr string) *http.ServeMux {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := db.New(dbConn)
	cfg := Config{
		db:        dbQueries,
		cache:     utils.NewLRUCache(CACHE_CAPACITY),
		cacheHits: 0,
		visitors: visitors{
			entries: make(map[string]*visitor),
		},
	}

	router := http.NewServeMux()

	router.Handle("POST /",
		cfg.rateLimiterMiddleware(http.HandlerFunc(cfg.handleShortening)))
	router.Handle("GET /{id}",
		cfg.rateLimiterMiddleware(http.HandlerFunc(cfg.handleRedirection)))

	go cfg.visitors.cleanupVisitors()
	return router
}
