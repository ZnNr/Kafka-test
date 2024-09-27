package database

import (
	"database/sql"
	"errors"
	"github.com/ZnNr/WB-test-L0/internal/models"
	"log"
)

// AddPayment добавляет информацию об оплате в базу данных.
func AddPayment(tx *sql.Tx, payment models.Payment, orderUID string) error {
	query := `INSERT INTO payments 
        ("transaction", "request_id", "currency", "provider", "amount", 
        "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee", "order_uid") 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// Выполнение SQL запроса и обработка ошибок.
	if _, err := tx.Exec(
		query,
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
	); err != nil {
		log.Printf("failed to add payment to database, Transaction:  %s, error: %v", payment.Transaction, err)
	}
	log.Printf("Successfully added payment to the database, transaction: %s", payment.Transaction)
	return nil
}

// GetPayment извлекает информацию об оплате из базы данных по UID заказа.
func GetPayment(db *sql.DB, orderUID string) (*models.Payment, error) {
	query := "SELECT * FROM payments WHERE order_uid = $1"
	row := db.QueryRow(query, orderUID)
	var payment models.Payment
	var uid string

	// Извлечение данных из строки результата и обработка ошибок.
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
			log.Printf("Payment not found for order UID: %s", orderUID)
			return nil, nil
		}
		log.Printf("Get payment failed for order UID: %s, error: %v", orderUID, err)
		return nil, err
	}
	return &payment, nil
}
