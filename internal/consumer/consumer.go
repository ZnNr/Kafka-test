package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ZnNr/WB-test-L0/internal/cache"
	"github.com/ZnNr/WB-test-L0/internal/models"
	"github.com/ZnNr/WB-test-L0/internal/repository"
	"go.uber.org/zap"
	"sync"
)

// Subscribe подписывается на сообщения Kafka и обрабатывает их.
func Subscribe(ctx context.Context, cache *cache.Cache, db *repository.OrdersRepo, logger *zap.Logger, wg *sync.WaitGroup) error {
	defer wg.Done() // Убедимся, что wait group завершится
	topic := "orders"

	worker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		return fmt.Errorf("failed to connect consumer: %w", err)
	}
	defer func() {
		if err := worker.Close(); err != nil {
			logger.Error("Failed to close consumer", zap.Error(err))
		}
	}()

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("consume_partition failed: %w", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Error("Failed to close partition consumer", zap.Error(err))
		}
	}()

	logger.Info("Consumer subscribed to Kafka!", zap.String("topic", topic))

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				logger.Error("Consuming error", zap.Error(err))
			case msg := <-consumer.Messages():
				handleMessage(msg, cache, db, logger)
			case <-ctx.Done():
				logger.Info("Shutting down consumer")
				return
			}
		}
	}()

	<-ctx.Done() // Ожидаем завершения по сигналу из контекста
	return nil
}

// handleMessage обрабатывает сообщение из Kafka.
func handleMessage(msg *sarama.ConsumerMessage, cache *cache.Cache, db *repository.OrdersRepo, logger *zap.Logger) {
	if len(msg.Value) == 0 {
		logger.Warn("Received empty message, skipping")
		return
	}

	var order models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		logger.Error("Failed to unmarshal message", zap.Error(err), zap.ByteString("message", msg.Value))
		return
	}

	if _, found := cache.GetOrder(order.OrderUID); found {
		logger.Info("Order exists, skipping", zap.String("order_uid", order.OrderUID))
		return
	}

	if err := db.AddOrder(order); err != nil {
		logger.Error("Failed to save order to DB", zap.Error(err), zap.String("order_uid", order.OrderUID))
		return
	}

	cache.SaveOrder(order)
	logger.Info("Consumed order", zap.String("order_uid", order.OrderUID))
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	return sarama.NewConsumer(brokers, config)
}
