package script

import (
	"fmt"
	"github.com/ZnNr/WB-test-L0/internal/models"
	"math/rand"
	"time"
)

func GenerateOrder() models.Order {
	delivery := models.Delivery{
		Name:    randomString(10),
		Phone:   randomPhone(),
		Zip:     randomZip(),
		City:    randomString(8),
		Address: randomString(15),
		Region:  randomString(8),
		Email:   randomString(5),
	}
	item := models.Item{
		ChrtID:      rand.Intn(1000),
		TrackNumber: randomString(10),
		Price:       float64(rand.Intn(1000)),
		Rid:         randomString(6),
		Name:        randomString(10),
		Sale:        float64(rand.Intn(100)),
		Size:        randomSize(),
		TotalPrice:  float64(rand.Intn(1000)),
		NmID:        rand.Intn(1000),
		Brand:       randomString(8),
		Status:      rand.Intn(5),
	}

	currencies := []string{"USD", "RUB", "EUR"}
	currency := currencies[rand.Intn(len(currencies))]

	payment := models.Payment{
		Transaction:  randomString(10),
		RequestID:    randomString(8),
		Currency:     currency,
		Provider:     randomString(6),
		Amount:       float64(rand.Intn(10000)),
		PaymentDT:    int(time.Now().Unix()),
		Bank:         randomString(6),
		DeliveryCost: float64(rand.Intn(500)),
		GoodsTotal:   float64(rand.Intn(10000)),
		CustomFee:    float64(rand.Intn(100)),
	}

	localies := []string{"en", "ru"}
	locale := localies[rand.Intn(len(localies))]
	order := models.Order{
		OrderUID:          randomString(12),
		TrackNumber:       randomString(10),
		Entry:             randomString(5),
		Delivery:          delivery,
		Payment:           payment,
		Items:             []models.Item{item},
		Locale:            locale,
		InternalSignature: randomString(8),
		CustomerID:        randomString(8),
		DeliveryService:   randomString(5),
		Shardkey:          randomString(5),
		SmID:              rand.Intn(100),
		DateCreated:       time.Now().Format("2006-01-02"),
		OofShard:          randomString(4),
	}
	return order
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Функция для генерации случайного номера телефона
func randomPhone() string {
	return fmt.Sprintf("+1%010d", rand.Int63n(10000000000))
}

// Функция для генерации случайного ZIP-кода
func randomZip() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

// Функция для генерации случайного размера
func randomSize() string {
	sizes := []string{"S", "M", "L", "XL", "XXL"}
	return sizes[rand.Intn(len(sizes))]
}
