package server

import (
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/cache"
	"github.com/ZnNr/WB-test-L0/internal/controller/router"
	"github.com/ZnNr/WB-test-L0/internal/repository/config"
	"log"
	"net/http"
)

type Server struct {
	cfg      config.ConfigApp
	Cache    *cache.Cache
	HTTPPort string
}

func New(cfgPath string, cache *cache.Cache) (*Server, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Server{
		cfg:      cfg.App,
		Cache:    cache,
		HTTPPort: fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
	}, nil
}

func (s *Server) Launch() error {
	r := router.NewController(s.Cache).SetupRouter()
	log.Printf("Starting server at %s\n", s.HTTPPort)

	err := http.ListenAndServe(s.HTTPPort, r)
	if err != nil {
		return fmt.Errorf("failed to launch server: %w", err)
	}
	return nil
}
