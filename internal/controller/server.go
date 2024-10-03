package controller

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/repository/config"
	"log"
	"net/http"

	"github.com/ZnNr/WB-test-L0/internal/cache"
	"github.com/gorilla/mux"
)

type Server struct {
	cfg   config.ConfigApp
	Cache *cache.Cache
}

func New(cfgPath string, cache *cache.Cache) (*Server, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &Server{
		cfg:   cfg.App,
		Cache: cache,
	}, nil
}

func (s *Server) Launch() error {
	r := s.setupRouter()
	http.Handle("/", r)
	log.Printf("Starting controller at %s:%s\n", s.cfg.Host, s.cfg.Port)

	err := http.ListenAndServe(s.cfg.Host+":"+s.cfg.Port, nil)
	if err != nil {
		return fmt.Errorf("failed to launch controller: %w", err)
	}
	return nil
}

func (s *Server) setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/order/{order_uid}", s.HandleGetOrder).Methods(http.MethodGet)
	return r
}

func (s *Server) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]
	order, ok := s.Cache.GetOrder(orderUID)
	if !ok {
		http.Error(w, fmt.Sprintf("OrderUID: <%s> not found!", orderUID), http.StatusNotFound)
		return
	}

	orderJSON, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		http.Error(w, "Failed to encode order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(orderJSON)
}

// хендлер для удаления заказа только из кэша
func (s *Server) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if !s.Cache.OrderExists(orderUID) {
		http.Error(w, fmt.Sprintf("OrderUID: <%s> not found!", orderUID), http.StatusNotFound)
		return
	}

	// Удаление заказа из кэша
	s.Cache.RemoveOrder(orderUID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("OrderUID: <%s> successfully deleted", orderUID)))
}

// Хендлер для очистки всех заказов из кэша
func (s *Server) HandleClearOrders(w http.ResponseWriter, r *http.Request) {
	s.Cache.Clear()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All orders successfully cleared"))
}
