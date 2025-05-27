package server


import (
	"net"
	"net/http"
	"sync"
	"time"
	"fmt"
)

const (
	rateLimit       = 60
	rateLimitWindow = time.Minute
	cleanupInterval = 5 * time.Minute
	expiryDuration  = 10 * time.Minute
)

type visitor struct {
	lastSeen time.Time
	mu sync.Mutex
	tokens int
}

type visitors struct {
	entries map[string]*visitor
	mu sync.Mutex
}

func (this *visitors) cleanupVisitors() {
	for {
		time.Sleep(cleanupInterval)

		this.mu.Lock()
		for ip, v := range this.entries {
			v.mu.Lock()
			if time.Since(v.lastSeen) > expiryDuration {
				delete(this.entries, ip)
			}
			v.mu.Unlock()
		}
		this.mu.Unlock()
	}
}

func (this *visitors) getVisitor(ip string) *visitor {
	this.mu.Lock()
	defer this.mu.Unlock()

	v, ok := this.entries[ip]
	if ok {
		if time.Since(v.lastSeen) > rateLimitWindow {
			v.lastSeen = time.Now()
			v.tokens = rateLimit
		}
		return v
	}
	this.entries[ip] = &visitor{
		lastSeen: time.Now(),
		tokens: rateLimit,
	}
	return this.entries[ip]
}

func (cfg *Config) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		v := cfg.visitors.getVisitor(ip)
		v.mu.Lock()
		defer v.mu.Unlock()

		if v.tokens > 0 {
			v.tokens--
			fmt.Println(ip, v.tokens)
			v.lastSeen = time.Now()
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Rate limit exceeded. Try again later.",
				http.StatusTooManyRequests)
		}
	})
}
