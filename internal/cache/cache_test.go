package cache

import (
	"github.com/ZnNr/WB-test-L0/internal/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// SaveOrder handles an order with an empty UID gracefully
func TestSaveOrderHandlesEmptyUID(t *testing.T) {
	// Arrange
	cache := New(10)
	order := models.Order{OrderUID: ""}

	// Act
	cache.SaveOrder(order)

	// Assert
	_, exists := cache.GetOrder("")
	if !exists {
		t.Errorf("Expected order with empty UID to be saved, but it was not found")
	}
}

// SaveOrder processes an order with special characters in the UID
func TestSaveOrderWithSpecialCharactersInUID(t *testing.T) {
	// Arrange
	cache := New(10)
	order := models.Order{OrderUID: "!@#$%^&*()"}

	// Act
	cache.SaveOrder(order)

	// Assert
	savedOrder, exists := cache.GetOrder("!@#$%^&*()")
	if !exists || savedOrder.OrderUID != "!@#$%^&*()" {
		t.Errorf("Expected order with special characters in UID to be saved, but it was not found")
	}
}

// SaveOrder handles concurrent writes without data races
func TestSaveOrderConcurrentWrites(t *testing.T) {
	// Arrange
	cache := New(10)
	order := models.Order{OrderUID: "test-order-1"}

	var wg sync.WaitGroup
	numRoutines := 100

	// Act
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.SaveOrder(order)
		}()
	}

	wg.Wait()

	// Assert
	savedOrder, exists := cache.GetOrder("test-order-1")
	if !exists || savedOrder.OrderUID != "test-order-1" {
		t.Errorf("Expected order with UID 'test-order-1' to be saved, but it was not found")
	}
}

// Retrieve an existing order successfully using a valid UID
func TestRetrieveExistingOrderWithValidUID(t *testing.T) {
	// Arrange
	cache := New(10)
	order := models.Order{OrderUID: "12345"}
	cache.SaveOrder(order)

	// Act
	retrievedOrder, exists := cache.GetOrder("12345")

	// Assert
	if !exists {
		t.Errorf("Expected order to exist, but it does not")
	}
	if retrievedOrder.OrderUID != "12345" {
		t.Errorf("Expected OrderUID to be '12345', got %s", retrievedOrder.OrderUID)
	}
}

// Attempt to retrieve an order with an empty UID string
func TestRetrieveOrderWithEmptyUID(t *testing.T) {
	// Arrange
	cache := New(10)

	// Act
	_, exists := cache.GetOrder("")

	// Assert
	if exists {
		t.Errorf("Expected order not to exist, but it does")
	}
}

// Retrieve an order when the cache is empty
func TestRetrieveOrderFromEmptyCache(t *testing.T) {
	// Arrange
	cache := New(10)

	// Act
	_, exists := cache.GetOrder("12345")

	// Assert
	if exists {
		t.Errorf("Expected order not to exist, but it does")
	}
}

// Returns true if the order exists in the cache
func TestOrderExistsReturnsTrue(t *testing.T) {
	// Arrange
	cache := New(10)
	order := models.Order{OrderUID: "12345"}
	cache.SaveOrder(order)

	// Act
	exists := cache.OrderExists("12345")

	// Assert
	if !exists {
		t.Errorf("Expected order to exist, but it does not")
	}
}

// Handles empty string as orderUID gracefully
func TestOrderExistsHandlesEmptyString(t *testing.T) {
	// Arrange
	cache := New(10)

	// Act
	exists := cache.OrderExists("")

	// Assert
	if exists {
		t.Errorf("Expected order not to exist for empty UID, but it does")
	}
}

// Handles very large orderUID strings without performance degradation
func TestOrderExistsHandlesLargeOrderUID(t *testing.T) {
	// Arrange
	cache := New(10)
	largeUID := string(make([]byte, 10000)) // Create a large string of 10,000 characters

	// Act
	exists := cache.OrderExists(largeUID)

	// Assert
	if exists {
		t.Errorf("Expected order not to exist for large UID, but it does")
	}
}

func TestClearEmptiesCache(t *testing.T) {
	// Arrange
	c := New(10)
	order := models.Order{OrderUID: "123"}
	c.SaveOrder(order)

	// Act
	c.Clear()

	// Assert
	_, exists := c.GetOrder("123")
	assert.False(t, exists, "Cache should be empty after Clear is called")
}

func TestClearConcurrent(t *testing.T) {
	// Arrange
	c := New(10)
	var wg sync.WaitGroup
	wg.Add(2)

	// Act
	go func() {
		defer wg.Done()
		c.Clear()
	}()

	go func() {
		defer wg.Done()
		c.Clear()
	}()

	wg.Wait()

	// Assert
	assert.Equal(t, 0, len(c.Orders), "Cache should be empty after concurrent Clear calls")
}

func TestClearDuringRead(t *testing.T) {
	// Arrange
	c := New(10)
	order := models.Order{OrderUID: "123"}
	c.SaveOrder(order)

	var wg sync.WaitGroup
	wg.Add(2)

	// Act
	go func() {
		defer wg.Done()
		c.Clear()
	}()

	go func() {
		defer wg.Done()
		_, exists := c.GetOrder("123")
		assert.False(t, exists, "Order should not exist after Clear is called")
	}()

	wg.Wait()
}
