package middleware

import (
	"github.com/mr-isik/gatling-backend/internal/api/httputil"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	tokens     int
	lastRefill time.Time
}

var (
	mu       sync.Mutex
	visitors = make(map[string]*visitor)
)

func RateLimitMiddleware(requestsPerSecond int) func(http.Handler) http.Handler {
	// Simple token bucket implementation
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				visitors[ip] = &visitor{
					tokens:     requestsPerSecond,
					lastRefill: time.Now(),
				}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			now := time.Now()
			elapsed := now.Sub(v.lastRefill).Seconds()

			// Refill tokens
			v.tokens += int(elapsed * float64(requestsPerSecond))
			if v.tokens > requestsPerSecond {
				v.tokens = requestsPerSecond
			}
			v.lastRefill = now

			if v.tokens > 0 {
				v.tokens--
				mu.Unlock()
				next.ServeHTTP(w, r)
			} else {
				mu.Unlock()
				httputil.JSONError(w, http.StatusTooManyRequests, domain.ErrBadRequest) // using bad request as generic, ideally would have ErrTooManyRequests
			}
		})
	}
}
