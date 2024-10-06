package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
)

const (
	getPaymentQuery = "SELECT \"transaction\", \"request_id\", \"currency\", \"provider\", \"amount\", \"payment_dt\", \"bank\", \"delivery_cost\", \"goods_total\", \"custom_fee\" FROM payments WHERE order_uid = $1"
	addPaymentQuery = `INSERT INTO payments
        ("transaction", "request_id", "currency", "provider", "amount",
        "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee", "order_uid")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
)

// AddPayment добавляет платеж в базу данных.
func AddPayment(tx *sql.Tx, payment models.Payment, orderUID string) error {
	_, err := tx.Exec(
		addPaymentQuery,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
		orderUID,
	)
	if err != nil {
		return fmt.Errorf("не удалось добавить платеж: %w", err) // Улучшено сообщение об ошибке
	}
	return nil
}

// GetPayment получает платеж из базы данных по orderUID.
func GetPayment(tx *sql.Tx, orderUID string) (*models.Payment, error) {
	row := tx.QueryRow(getPaymentQuery, orderUID) // Используем tx
	var payment models.Payment

	err := row.Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("платеж не найден: %w", err)
		}
		return nil, fmt.Errorf("не удалось получить платеж: %w", err)
	}
	return &payment, nil
}

// PaymentExists проверяет существование платежа в базе данных по orderUID.
func PaymentExists(tx *sql.Tx, orderUID string) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM payments WHERE order_uid = $1)", orderUID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить существование платежа: %w", err) // Улучшено сообщение об ошибке
	}
	return exists, nil
}
