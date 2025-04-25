package main

import (
	"log/slog"
	"net/http"
)

func httpLogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rw := responseWriter{
			ResponseWriter: res,
		}
		h.ServeHTTP(&rw, req)
		slog.Info("http request",
			"method", req.Method,
			"url", req.URL.String(),
			"status", rw.status,
			"status_text", http.StatusText(rw.status),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(p)
}
