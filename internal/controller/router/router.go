package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"

	"github.com/ZnNr/WB-test-L0/internal/cache"
	"github.com/gorilla/mux"
)

type Controller struct {
	Cache *cache.Cache
}

func NewController(cache *cache.Cache) *Controller {
	return &Controller{Cache: cache}
}

func (c *Controller) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Настройка CORS
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	// Применяем middleware для CORS
	r.Use(handlers.CORS(corsOptions, corsMethods, corsHeaders))

	// Маршруты вашего API
	r.HandleFunc("/order/{order_uid}", c.HandleGetOrder).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/order/{order_uid}", c.HandleDeleteOrder).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/orders", c.HandleClearOrders).Methods(http.MethodDelete, http.MethodOptions)

	return r
}

func (c *Controller) HandleGetOrder(w http.ResponseWriter, r *http.Request) {

	orderUID := mux.Vars(r)["order_uid"]

	order, ok := c.Cache.GetOrder(orderUID)
	if !ok {
		c.writeError(w, http.StatusNotFound, fmt.Sprintf("OrderUID: <%s> not found!", orderUID))
		return
	}
	w.Write([]byte("Order details for UID: " + orderUID))
	c.writeJSON(w, http.StatusOK, order)
}

func (c *Controller) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Просто ответить на preflight запрос
		return
	}

	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if !c.Cache.OrderExists(orderUID) {
		c.writeError(w, http.StatusNotFound, fmt.Sprintf("OrderUID: <%s> not found!", orderUID))
		return
	}

	c.Cache.RemoveOrder(orderUID)
	c.writeJSON(w, http.StatusOK, fmt.Sprintf("OrderUID: <%s> successfully deleted", orderUID))
}

func (c *Controller) HandleClearOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Просто ответить на preflight запрос
		return
	}

	c.Cache.Clear()
	c.writeJSON(w, http.StatusOK, "All orders successfully cleared")
}

func (c *Controller) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) writeError(w http.ResponseWriter, status int, message string) {
	c.writeJSON(w, status, map[string]string{"error": message})
}
