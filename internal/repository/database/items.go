package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
)

const (
	addItemQuery     = `INSERT INTO items ("chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status", "order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	getAllItemsQuery = "SELECT * FROM items WHERE order_uid = $1"
)

func AddItems(tx *sql.Tx, items []models.Item, OrderUID string) (err error) {
	for _, item := range items {
		err = AddItem(tx, item, OrderUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddItem(tx *sql.Tx, item models.Item, OrderUID string) error {
	//query := `INSERT INTO "items"("chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status", "order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := tx.Exec(
		addItemQuery,
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
		OrderUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetItems(db *sql.DB, OrderUID string) ([]models.Item, error) {

	rows, err := db.Query(getAllItemsQuery, OrderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("items not found: %w", err)
		}
		return nil, fmt.Errorf("get items failed: %w", err)
	}
	var items []models.Item
	var uid string
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&uid,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
