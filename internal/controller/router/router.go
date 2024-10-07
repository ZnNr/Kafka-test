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

// Функция для инициализации контроллера с кэшем
func NewController(cache *cache.Cache) *Controller {
	return &Controller{Cache: cache}
}

// Настройка маршрутизатора
func (c *Controller) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Настройка CORS
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	//Применяем middleware для CORS
	r.Use(handlers.CORS(corsOptions, corsMethods, corsHeaders))
	r.Use(c.preflightHandler)

	// Маршруты вашего API
	r.HandleFunc("/order/{order_uid}", c.HandleGetOrder).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/order/{order_uid}", c.HandleDeleteOrder).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/delorders", c.HandleClearOrders).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/orders", c.HandleGetAllOrders).Methods(http.MethodGet, http.MethodOptions)

	return r
}

// Middleware для обработки предварительных запросов
func (c *Controller) preflightHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // Возвращаем статус 204 No Content
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HandleGetOrder Обработчик для получения заказа по order_uid
func (c *Controller) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := mux.Vars(r)["order_uid"]

	order, ok := c.Cache.GetOrder(orderUID)
	if !ok {
		c.writeError(w, http.StatusNotFound, fmt.Sprintf("OrderUID: <%s> not found!", orderUID))
		return
	}
	c.writeJSON(w, http.StatusOK, order)
}

// HandleDeleteOrder Обработчик для удаления заказа по order_uid
func (c *Controller) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if !c.Cache.OrderExists(orderUID) {
		c.writeError(w, http.StatusNotFound, fmt.Sprintf("OrderUID: <%s> not found!", orderUID))
		return
	}

	c.Cache.RemoveOrder(orderUID)
	c.writeJSON(w, http.StatusOK, fmt.Sprintf("OrderUID: <%s> successfully deleted", orderUID))
}

// HandleClearOrders Обработчик для очистки всех заказов
func (c *Controller) HandleClearOrders(w http.ResponseWriter, r *http.Request) {
	c.Cache.Clear()
	c.writeJSON(w, http.StatusOK, "All orders successfully cleared")
}

// HandleGetAllOrders обработчик для получения всех заказов
func (c *Controller) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := c.Cache.GetAllOrders() // Метод, который должен вернуть все заказы
	if len(orders) == 0 {
		c.writeJSON(w, http.StatusOK, []interface{}{}) // Если заказов нет, возвращаем пустой массив
		return
	}
	c.writeJSON(w, http.StatusOK, orders)
}

// Приватные методы для записи JSON и ошибок
func (c *Controller) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) writeError(w http.ResponseWriter, status int, message string) {
	c.writeJSON(w, status, map[string]string{"error": message})
}
