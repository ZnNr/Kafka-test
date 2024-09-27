package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
	"github.com/ZnNr/WB-test-L0/internal/repository/config"
	"github.com/ZnNr/WB-test-L0/internal/repository/database"
)

type OrdersRepo struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*OrdersRepo, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return &OrdersRepo{DB: db}, nil
}
func (o *OrdersRepo) AddOrder(order models.Order) error {
	tx, err := o.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(
		query,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := database.AddPayment(tx, order.Payment, order.OrderUID); err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	if err := database.AddItems(tx, order.Items, order.OrderUID); err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}

	if err := database.AddDelivery(tx, order.Delivery, order.OrderUID); err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	return tx.Commit()
}

func (o *OrdersRepo) GetOrder(orderUID string) (*models.Order, error) {
	query := "SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1"
	row := o.DB.QueryRow(query, orderUID)

	var order models.Order
	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Delivery, err = database.GetDelivery(o.DB, orderUID); err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	if order.Payment, err = database.GetPayment(o.DB, orderUID); err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if order.Items, err = database.GetItems(o.DB, orderUID); err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return &order, nil
}
func (o *OrdersRepo) GetOrders() ([]models.Order, error) {
	query := "SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders"
	rows, err := o.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard); err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		if order.Delivery, err = database.GetDelivery(o.DB, order.OrderUID); err != nil {
			return nil, fmt.Errorf("failed to get delivery for order %s: %w", order.OrderUID, err)
		}

		if order.Payment, err = database.GetPayment(o.DB, order.OrderUID); err != nil {
			return nil, fmt.Errorf("failed to get payment for order %s: %w", order.OrderUID, err)
		}

		if order.Items, err = database.GetItems(o.DB, order.OrderUID); err != nil {
			return nil, fmt.Errorf("failed to get items for order %s: %w", order.OrderUID, err)
		}

		orders = append(orders, order)
	}
	return orders, nil
}
