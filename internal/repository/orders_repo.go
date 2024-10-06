package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
	"github.com/ZnNr/WB-test-L0/internal/repository/config"
	"github.com/ZnNr/WB-test-L0/internal/repository/database"
	_ "github.com/lib/pq"
)

const (
	addOrderQuery     = `INSERT INTO orders("order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	getOrderQuery     = "SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1"
	getAllOrdersQuery = "SELECT * FROM orders"
)

type OrdersRepo struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*OrdersRepo, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &OrdersRepo{DB: db}, nil
}

func (o *OrdersRepo) OrderExists(orderUID string) (bool, error) {
	var exists bool
	err := o.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", orderUID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (o *OrdersRepo) AddOrder(order models.Order) error {
	exists, err := o.OrderExists(order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to check if order exists: %w", err)
	}

	if exists {
		return fmt.Errorf("order with order_uid %s already exists", order.OrderUID)
	}

	tx, err := o.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if commit does not happen.

	_, err = tx.Exec(addOrderQuery, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Payment existence check and insertion.
	if err := o.processPayment(tx, order); err != nil {
		return err
	}

	if err := database.AddItems(tx, order.Items, order.OrderUID); err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}

	if err := database.AddDelivery(tx, order.Delivery, order.OrderUID); err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// processPayment checks if the payment exists and adds a new payment record if it doesn't.
func (o *OrdersRepo) processPayment(tx *sql.Tx, order models.Order) error {
	exists, err := database.PaymentExists(tx, order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to check if payment exists: %w", err)
	}

	if !exists {
		if err := database.AddPayment(tx, order.Payment, order.OrderUID); err != nil {
			return fmt.Errorf("failed to insert payment: %w", err)
		}
	}
	// Future logic for updating an existing payment can be added here.

	return nil
}
func (o *OrdersRepo) GetOrder(orderUID string) (*models.Order, error) {
	var order models.Order
	tx, err := o.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := tx.QueryRow(getOrderQuery, orderUID).Scan(&order.OrderUID,
		&order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID,
		&order.DateCreated, &order.OofShard); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := populateOrderDetails(tx, &order); err != nil { // Передайте транзакцию
		return nil, fmt.Errorf("failed to populate order details: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}
func populateOrderDetails(tx *sql.Tx, order *models.Order) error {
	delivery, err := database.GetDelivery(tx, order.OrderUID)
	if err != nil {
		return err
	}
	order.Delivery = *delivery

	payment, err := database.GetPayment(tx, order.OrderUID) // Уже правильно
	if err != nil {
		return err
	}
	order.Payment = *payment

	items, err := database.GetItems(tx, order.OrderUID) // Уже правильно
	if err != nil {
		return err
	}
	order.Items = items

	return nil
}

func (o *OrdersRepo) GetOrders() ([]models.Order, error) {
	tx, err := o.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(getAllOrdersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey,
			&order.SmID, &order.DateCreated, &order.OofShard); err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		if err := populateOrderDetails(tx, &order); err != nil { // Передайте транзакцию
			return nil, fmt.Errorf("failed to get details for order %s: %w", order.OrderUID, err)
		}

		orders = append(orders, order)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return orders, nil
}
