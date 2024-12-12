package controllers

import (
	"MiniHIFPT/models"
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

func GetCustomers(c *fiber.Ctx) error {
	customers, err := services.GetAllCustomers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin khách hàng",
		})
	}
	return c.JSON(customers)
}

func CreateCustomers(c *fiber.Ctx) error {
	var tempCustomer models.TempCustomer

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&tempCustomer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	customer, err := services.CreateCustomerService(&tempCustomer)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm khách hàng thành công",
		"data":    customer,
	})
}
