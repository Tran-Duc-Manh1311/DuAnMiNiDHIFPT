package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

type ContractController struct {
	Service *services.ContractService
}

// Thống kê hợp đồng theo trạng thái
func (controller *ContractController) GetContractCountByStatus(c *fiber.Ctx) error {
	status := c.Params("status") // Lấy trạng thái từ URL
	count, err := controller.Service.CountContractsByStatus(status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
		"count":  count,
	})
}
