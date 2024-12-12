package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// API chuyển sở hữu hợp đồng
func TransferOwnership(c *fiber.Ctx) error {
	// Lấy tham số từ yêu cầu JSON
	var request struct {
		OldCustomerID string `json:"oldCustomerId"`
		NewCustomerID string `json:"newCustomerId"`
	}

	// Lấy accountID từ thông tin người dùng đã đăng nhập (tương tự như trong UpdateContract)
	accountID := c.Locals("accountID").(string)

	// Phân tích dữ liệu JSON từ yêu cầu
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Giá trị nhập vào không hợp lệ",
		})
	}

	// Gọi service để xử lý logic chuyển sở hữu hợp đồng
	err := services.TransferOwnership(accountID, request.OldCustomerID, request.NewCustomerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Trả về kết quả
	return c.JSON(fiber.Map{
		"message": "Tất cả hợp đồng đã được chuyển giao thành công",
	})
}
