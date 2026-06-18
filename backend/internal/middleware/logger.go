package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLog logs every incoming request with method, path, status and duration.
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wr, r)
		log.Printf("[%s] %s %s → %d (%v)", r.Method, r.URL.Path, r.RemoteAddr, wr.status, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
