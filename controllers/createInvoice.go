package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Tạo hóa đơn mới
func CreateInvoice(c *fiber.Ctx) error {
	var invoiceRequest struct {
		ContractID  string  `json:"contract_id"`
		Amount      float64 `json:"amount"`
		ServiceName string  `json:"service_name"`
	}

	// Parse dữ liệu đầu vào
	if err := c.BodyParser(&invoiceRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Tạo ngày đến hạn mặc định (30 ngày kể từ ngày hiện tại)
	dueDate := time.Now().Add(30 * 24 * time.Hour) // Tính ngày đến hạn là 30 ngày sau ngày hiện tại

	// Tạo hóa đơn
	invoice, err := services.CreateInvoice(invoiceRequest.ContractID, invoiceRequest.Amount, dueDate, invoiceRequest.ServiceName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo hóa đơn",
		})
	}

	// Trả về kết quả
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tạo hóa đơn thành công",
		"invoice": invoice,
	})
}
