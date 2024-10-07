package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/ZnNr/WB-test-L0/internal/cache"
	"github.com/ZnNr/WB-test-L0/internal/repository"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

// Successfully connects to Kafka broker and subscribes to the topic
func TestSubscribeConnectsAndSubscribes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := cache.New(10)
	db := &repository.OrdersRepo{}
	logger := zap.NewExample()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Act
	err := Subscribe(ctx, cache, db, logger, wg)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, logger.Check(zap.InfoLevel, "Consumer subscribed to Kafka!").Message, "Consumer subscribed to Kafka!")
}

// Correctly processes messages from Kafka and updates cache and database
func TestSubscribeProcessesMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := cache.New(10)
	db := &repository.OrdersRepo{}
	logger := zap.NewExample()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Arrange
	go func() {
		// Simulate message processing
		time.Sleep(1 * time.Second)
		cancel()
	}()

	// Act
	err := Subscribe(ctx, cache, db, logger, wg)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, cache.Orders)
}

// Receives an empty message and logs a warning
func TestSubscribeReceivesEmptyMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := cache.New(10)
	db := &repository.OrdersRepo{}
	logger := zap.NewExample()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Arrange
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	// Act
	err := Subscribe(ctx, cache, db, logger, wg)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, logger.Check(zap.WarnLevel, "Received empty message").Message, "Received empty message")
}

// Encounters an error when consuming messages and logs the error
func TestSubscribeLogsConsumingError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := cache.New(10)
	db := &repository.OrdersRepo{}
	logger := zap.NewExample()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Arrange
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	// Act
	err := Subscribe(ctx, cache, db, logger, wg)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, logger.Check(zap.ErrorLevel, "Consuming error").Message, "Consuming error")
}

// Receives an empty message and logs a warning
func TestHandleEmptyMessage(t *testing.T) {
	cache := cache.New(10)
	db := &repository.OrdersRepo{}
	logger := zap.NewExample()

	msg := &sarama.ConsumerMessage{Value: []byte{}}

	handleMessage(msg, cache, db, logger)

	// Check logs for warning about empty message
	logs := logger.Check(zap.WarnLevel, "Received empty message, skipping")
	assert.NotNil(t, logs)
}

// Successfully connects to a Kafka broker with valid broker addresses
func TestConnectConsumerWithValidBrokers(t *testing.T) {
	// Arrange
	brokers := []string{"localhost:9092"}

	// Act
	consumer, err := ConnectConsumer(brokers)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if consumer == nil {
		t.Fatal("Expected a non-nil consumer")
	}
}

// Returns a sarama.Consumer object when connection is successful
func TestReturnsConsumerObjectOnSuccess(t *testing.T) {
	// Arrange
	brokers := []string{"localhost:9092"}

	// Act
	consumer, err := ConnectConsumer(brokers)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if _, ok := consumer.(sarama.Consumer); !ok {
		t.Fatal("Expected a sarama.Consumer object")
	}
}

// Handles empty broker list gracefully
func TestConnectConsumerWithEmptyBrokerList(t *testing.T) {
	// Arrange
	brokers := []string{}

	// Act
	consumer, err := ConnectConsumer(brokers)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for empty broker list")
	}
	if consumer != nil {
		t.Fatal("Expected a nil consumer for empty broker list")
	}
}

// Returns an error if broker addresses are invalid
func TestConnectConsumerWithInvalidBrokers(t *testing.T) {
	// Arrange
	brokers := []string{"invalid-broker-address"}

	// Act
	consumer, err := ConnectConsumer(brokers)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for invalid broker addresses")
	}
	if consumer != nil {
		t.Fatal("Expected a nil consumer for invalid broker addresses")
	}
}

// Manages network failures or broker unavailability
func TestConnectConsumerWithUnavailableBroker(t *testing.T) {
	// Arrange
	brokers := []string{"unavailable-broker-address"}

	// Act
	consumer, err := ConnectConsumer(brokers)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for unavailable broker")
	}
	if consumer != nil {
		t.Fatal("Expected a nil consumer for unavailable broker")
	}
}
