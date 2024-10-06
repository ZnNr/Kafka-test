package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZnNr/WB-test-L0/internal/models"
)

const (
	addDeliveryQuery = `INSERT INTO deliveries
    ("name", "phone", "zip", "city", "address", "region", "email", "order_uid")
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (order_uid) DO UPDATE SET
        name = EXCLUDED.name,
        phone = EXCLUDED.phone,
        zip = EXCLUDED.zip,
        city = EXCLUDED.city,
        address = EXCLUDED.address,
        region = EXCLUDED.region,
        email = EXCLUDED.email
    RETURNING order_uid` // Возвращаем order_uid после вставки/обновления

	getDeliveryQuery = `SELECT order_uid, name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid = $1`
)

func AddDelivery(tx *sql.Tx, delivery models.Delivery, orderUID string) error {
	// Проверяем, существует ли доставка с данным order_uid
	existingDelivery, err := GetDelivery(tx, orderUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get delivery: %w", err)
	}

	if existingDelivery != nil {
		return fmt.Errorf("delivery with order_uid %s already exists", orderUID)
	}

	// Вставляем новую доставку или обновляем
	var returnedOrderUID string
	err = tx.QueryRow(
		addDeliveryQuery,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		orderUID,
	).Scan(&returnedOrderUID)

	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	fmt.Printf("Inserted/Updated delivery with order_uid: %s\n", returnedOrderUID)
	return nil
}

func GetDelivery(tx *sql.Tx, orderUID string) (*models.Delivery, error) {
	row := tx.QueryRow(getDeliveryQuery, orderUID)

	var delivery models.Delivery
	err := row.Scan(
		&delivery.OrderUID, // Убедитесь, что OrderUID присутствует в модели
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	fmt.Printf("Executing query: %s with orderUID: %s\n", getDeliveryQuery, orderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("get delivery failed: %w", err)
	}

	fmt.Printf("Scanned delivery: %+v\n", delivery)
	return &delivery, nil
}
