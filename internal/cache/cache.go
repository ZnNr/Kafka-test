package cache

import (
	"github.com/ZnNr/WB-test-L0/internal/models"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func New() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

// SaveOrder сохраняет заказ в кэш
func (c *Cache) SaveOrder(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

// GetOrder получает заказ из кэша по UID
func (c *Cache) GetdOrder(OrderUid string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[OrderUid]
	return order, ok
}
