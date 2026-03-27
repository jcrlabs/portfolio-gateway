package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
	"github.com/jonathanCaamano/portfolio-gateway/internal/gateway"
	"github.com/jonathanCaamano/portfolio-gateway/internal/metrics"
)

func main() {
	cfg := config.Load()
	routes := config.LoadRoutes(cfg.RoutingFile)

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	m := metrics.New()
	mux := gateway.NewMux(routes, cfg, log, m)

	log.Info("gateway starting", "port", cfg.Port, "routes", len(routes))
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
