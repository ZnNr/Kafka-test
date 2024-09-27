package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
)

// AddDelivery добавляет запись о доставке в базу данных.
func AddDelivery(tx *sql.Tx, delivery models.Delivery, orderUID string) error {
	query := `INSERT INTO deliveries
        ("name", "phone", "zip", "city", "address", "region", "email", "order_uid")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	if _, err := tx.Exec(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email, orderUID); err != nil {
		return fmt.Errorf("failed to add delivery: %w", err)
	}
	return nil
}

// GetDelivery извлекает данные о доставке из базы данных.
func GetDelivery(db *sql.DB, orderUID string) (*models.Delivery, error) {
	query := "SELECT * FROM deliveries WHERE order_uid = $1"

	row := db.QueryRow(query, orderUID)

	var delivery models.Delivery
	if err := row.Scan(&orderUID, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no delivery found for order %s: %w", orderUID, err)
		}
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	return &delivery, nil
}
