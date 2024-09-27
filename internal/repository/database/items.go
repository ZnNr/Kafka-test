package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
)

// AddItems добавляет несколько элементов в базу данных.
func AddItems(tx *sql.Tx, items []models.Item, orderUID string) error {
	for _, item := range items {
		if err := addItem(tx, item, orderUID); err != nil {
			return fmt.Errorf("failed to add item: %w", err)
		}
	}
	return nil
}

// addItem добавляет один элемент в базу данных.
func addItem(tx *sql.Tx, item models.Item, orderUID string) error {
	query := `INSERT INTO items
        ("chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status", "order_uid")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if _, err := tx.Exec(query,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
		orderUID); err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// GetItems извлекает элементы из базы данных по идентификатору заказа.
func GetItems(db *sql.DB, orderUID string) ([]models.Item, error) {
	query := "SELECT * FROM items WHERE order_uid = $1"
	rows, err := db.Query(query, orderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("items not found: %w", err)
		}
		return nil, fmt.Errorf("query items failed: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&orderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	// Проверка на наличие ошибок завершения итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return items, nil
}
