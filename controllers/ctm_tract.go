package controllers

import (
	"MiniHIFPT/models"
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các hợp đồng
func GetCtmContracts(c *fiber.Ctx) error {
	ctmContracts, err := services.GetCtmContracts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin hợp đồng",
		})
	}
	return c.JSON(ctmContracts)
}

// Tạo liên kết hợp đồng
func CreateCtmContracts(c *fiber.Ctx) error {
	// Khởi tạo một biến Customer_Contractt
	var ctmContract models.Customer_Contractt

	// Phân tích dữ liệu đầu vào từ body request
	if err := c.BodyParser(&ctmContract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Lấy accountID từ thông tin người dùng đã đăng nhập (tương tự như trong UpdateContract)
	accountID := c.Locals("accountID").(string)
	// Kiểm tra và tạo liên kết hợp đồng qua service
	if err := services.CreateCtmContract(ctmContract, accountID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Trả về kết quả thành công
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo liên kết số điện thoại và hợp đồng thành công",
	})
}
