package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
func GetAllInvoiceByID(invoiceID string) *ServiceResponse {
	// Chuyển đổi customerID thành UUID
	idUUID, err := uuid.Parse(invoiceID)
	if err != nil {
		return respond(fiber.StatusBadRequest, "ID hóa đơn không hợp lệ", nil)
	}

	// Gọi hàm GetCustomerByID để lấy thông tin khách hàng
	invoice, err := database.GetInvoiceByID(idUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respond(fiber.StatusNotFound, "Hóa đơn không tồn tại", nil)
		}
		return respond(fiber.StatusInternalServerError, "Lỗi khi lấy thông tin hóa đơn", nil)
	}

	return respond(fiber.StatusOK, "Lấy thông tin hóa đơn thành công", invoice)
}
