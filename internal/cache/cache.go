package cache

import (
	"github.com/ZnNr/WB-test-L0/internal/models"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	Orders map[string]models.Order
}

// New создает новый кэш с возможностью задания начальной ёмкости
func New(initialCapacity int) *Cache {
	return &Cache{
		Orders: make(map[string]models.Order, initialCapacity),
	}
}

// SaveOrder сохраняет заказ в кэш
func (c *Cache) SaveOrder(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Orders[order.OrderUID] = order
}

// GetOrder получает заказ из кэша по UID
func (c *Cache) GetOrder(OrderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.Orders[OrderUID]
	return order, ok
}

func (c *Cache) OrderExists(orderUID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.Orders[orderUID]
	return exists
}

// RemoveOrder удаляет заказ из кэша по UID
func (c *Cache) RemoveOrder(OrderUID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Orders, OrderUID)
}

// Clear очищает кэш
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Orders = make(map[string]models.Order)
}
