package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// API Tìm kiếm khách hàng/hợp đồng
func SearchContractOrCustomer(c *fiber.Ctx) error {
	searchTerm := c.Query("searchTerm") // Lấy từ query parameter

	// Kiểm tra nếu searchTerm rỗng
	if searchTerm == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Vui lòng nhập từ khóa tìm kiếm",
		})
	}

	// Tìm kiếm khách hàng và hợp đồng qua service
	customers, contracts, err := services.SearchContractsAndCustomers(searchTerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Trả về kết quả tìm kiếm
	return c.JSON(fiber.Map{
		"customers": customers,
		"contracts": contracts,
	})
}
