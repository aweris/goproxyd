package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/goproxy/goproxy"
)

func main() {
	var (
		root      string
		addr      string
		readonly  bool
		logFormat string
	)
	flag.StringVar(&root,
		"root", "/cache/download",
		"cache download root",
	)
	flag.BoolVar(&readonly,
		"ro", true,
		"readonly mode",
	)
	flag.StringVar(&addr,
		"addr", "0.0.0.0:8080",
		"addr to bind to",
	)
	flag.StringVar(&logFormat,
		"log-format", "text",
		"log format (text or json)",
	)
	flag.Parse()

	// Configure logger
	var logHandler slog.Handler
	if logFormat == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	proxy := &goproxy.Goproxy{
		Cacher: &cacher{
			root:     root,
			readonly: readonly,
		},
		Logger: logger,
	}

	logger.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, httpLogHandler(proxy)); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
