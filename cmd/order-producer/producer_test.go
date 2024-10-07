package main

import (
	"github.com/IBM/sarama"

	"testing"
)

// Successfully connect to a Kafka broker using valid broker addresses
func TestConnectProducerWithValidBrokers(t *testing.T) {
	// Arrange
	brokers := []string{"localhost:9092"}

	// Act
	producer, err := ConnectProducer(brokers)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if producer == nil {
		t.Fatal("Expected a valid producer, got nil")
	}
}

// Return a SyncProducer instance when connection is successful
func TestReturnSyncProducerInstanceOnSuccess(t *testing.T) {
	// Arrange
	brokers := []string{"localhost:9092"}

	// Act
	producer, err := ConnectProducer(brokers)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if _, ok := producer.(sarama.SyncProducer); !ok {
		t.Fatal("Expected a SyncProducer instance")
	}
}

// Handle empty broker list gracefully
func TestHandleEmptyBrokerListGracefully(t *testing.T) {
	// Arrange
	brokers := []string{}

	// Act
	producer, err := ConnectProducer(brokers)

	// Assert
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
	if producer != nil {
		t.Fatal("Expected nil producer, got a valid instance")
	}
}

// Handle invalid broker addresses
func TestHandleInvalidBrokerAddresses(t *testing.T) {
	// Arrange
	brokers := []string{"invalid-broker"}

	// Act
	producer, err := ConnectProducer(brokers)

	// Assert
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
	if producer != nil {
		t.Fatal("Expected nil producer, got a valid instance")
	}
}
