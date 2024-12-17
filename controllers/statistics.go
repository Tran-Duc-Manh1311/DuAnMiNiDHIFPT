package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// API Thống kê hợp đồng theo trạng thái
func GetContractCountByStatus(c *fiber.Ctx) error {
	// Lấy trạng thái hợp đồng từ URL parameter
	status := c.Params("status")

	// Kiểm tra nếu status rỗng
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Vui lòng nhập trạng thái hợp đồng",
		})
	}

	// Gọi service để lấy số lượng hợp đồng theo trạng thái
	count, err := services.CountContractsByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Trả về kết quả thống kê
	return c.JSON(fiber.Map{
		"status": status,
		"count":  count,
	})
}
