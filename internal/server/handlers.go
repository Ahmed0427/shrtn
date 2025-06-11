package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/ahmed0427/shrtn/internal/db"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

const ID_DEFAULT_LENGTH = 8

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortID(length int) (string, error) {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", nil
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

func isValidURL(strURL string) bool {
	URL, err := url.Parse(strURL)
	if err != nil {
		return false
	}
	if URL.Scheme == "" || URL.Host == "" {
		return false
	}
	return true
}

func (cfg *Config) handleShortening(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "Missing 'url' field", http.StatusBadRequest)
		return
	}
	if !isValidURL(req.URL) {
		http.Error(w, "Not a valid URL", http.StatusBadRequest)
		return
	}

	var shortID string
	shortID, err := cfg.db.GetID(context.Background(), req.URL)
	if shortID == "" {
		for {
			shortID, err = generateShortID(ID_DEFAULT_LENGTH)
			if err != nil {
				http.Error(w, "Faild to generate short ID",
					http.StatusInternalServerError)
				return
			}
			original, _ := cfg.db.GetOriginalURL(context.Background(), shortID)
			if original == "" {
				params := db.AddURLParams{
					ID:             shortID,
					OriginalUrl:    req.URL,
					CreatedAt:      time.Now(),
					LastAccessedAt: time.Now(),
					AccessCount:    0,
				}
				_, err := cfg.db.AddURL(context.Background(), params)
				if err != nil {
					http.Error(w, "Faild to add entry to the database",
						http.StatusInternalServerError)
					return
				}
				cfg.cache.Put(shortID, req.URL)
				break
			}
		}
	}

	baseURL := "http://localhost:8080"
	shortURL := baseURL + "/" + shortID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": shortURL,
	})
}

func (cfg *Config) handleRedirection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	original_url, err := cfg.cache.Get(id)
	if err != nil {
		original_url, err = cfg.db.GetOriginalURL(context.Background(), id)
		if err != nil || original_url == "" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	} else {
		cfg.cacheHits += 1
	}
	cfg.cache.Put(id, original_url)
	cfg.db.UpdateAccessStats(context.Background(), id)
	w.Header().Set("Location", original_url)
	w.WriteHeader(http.StatusMovedPermanently)
}
