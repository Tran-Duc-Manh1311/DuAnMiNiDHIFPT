package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// Xử lý yêu cầu thanh toán
func ProcessPayment(c *fiber.Ctx) error {
	// Định nghĩa kiểu dữ liệu request
	var paymentRequest struct {
		InvoiceID string  `json:"invoice_id"`
		Amount    float64 `json:"amount"`
		Method    string  `json:"method"`
	}

	// Parse dữ liệu đầu vào từ body request
	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra tính hợp lệ của dữ liệu thanh toán
	if paymentRequest.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số tiền thanh toán phải lớn hơn 0",
		})
	}
	// Kiểm tra phương thức thanh toán
	if paymentRequest.Method == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phương thức thanh toán không hợp lệ",
		})
	}
	// Lấy accountID từ thông tin người dùng đã đăng nhập (tương tự như trong UpdateContract)
	accountID := c.Locals("accountID").(string)
	// Gọi service để xử lý thanh toán
	err := services.ProcessPayment(paymentRequest.InvoiceID, paymentRequest.Amount, paymentRequest.Method, accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể xử lý thanh toán: " + err.Error(),
		})
	}

	// Trả về kết quả thành công
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Thanh toán thành công",
	})
}
