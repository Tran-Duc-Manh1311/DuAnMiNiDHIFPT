package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// Tạo phương thức thanh toán
func CreatePaymentMethod(c *fiber.Ctx) error {
	err := services.CreateMethodService(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Message,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo phương thức thanh toán thành công",
	})
}
