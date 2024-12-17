package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Tạo hóa đơn mới
func CreateInvoice(contractID string, amount float64, dueDate time.Time, serviceName string) (*models.Invoice, error) {
	invoice := models.Invoice{
		ID:            uuid.New().String(),
		ContractID:    contractID,
		Amount:        amount,
		PaymentStatus: "Pending", // Trạng thái mặc định là Pending
		DueDate:       dueDate,
		ServiceName:   serviceName,
	}
	// Lưu hóa đơn vào cơ sở dữ liệu
	if err := database.DB.Create(&invoice).Error; err != nil {
		return nil, fmt.Errorf("không thể lưu hóa đơn vào cơ sở dữ liệu: %v", err)
	}

	return &invoice, nil
}
