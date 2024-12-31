package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
	// "regexp"
)

func CreateMethodService(c *fiber.Ctx) *ServiceResponse {
	var method models.PaymentMethod
	if err := c.BodyParser(&method); err != nil {
		return respond(fiber.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", nil)
	}

	if method.Method == "" || method.Name == "" {
		return respond(fiber.StatusBadRequest, "Thiếu thông tin cần thiết", nil)
	}

	existingmethod, err := database.Methoddetails(method.Name)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra phương thức thanh toán", nil)
	}
	if existingmethod != nil {
		return respond(fiber.StatusConflict, "Phương thức thanh toán đã tồn tại", nil)
	}

	if err := database.CreatePMMethod(&method); err != nil {
		return respond(fiber.StatusInternalServerError, "Không thể lưu phương thức thanh toán", nil)
	}

	return respond(fiber.StatusCreated, "Tạo phương thức thanh toán thành công", &method)
}
