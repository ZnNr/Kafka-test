package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
)

const (
	getPaymentQuery = "SELECT *  FROM payments WHERE order_uid = $1"
	addPaymentQuery = `INSERT INTO payments
("transaction", "request_id", "currency", "provider", "amount",
"payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee", "order_uid")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
)

func AddPayment(tx *sql.Tx, payment models.Payment, OrderUID string) error {

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
		OrderUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetPayment(db *sql.DB, OrderUID string) (*models.Payment, error) {

	row := db.QueryRow(getPaymentQuery, OrderUID)
	var payment models.Payment
	var uid string
	err := row.Scan(
		&uid,
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
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("get payment failed: %w", err)
	}
	return &payment, nil
}
