package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// AddContractAccessController xử lý thêm quyền truy cập cho tài khoản đối với hợp đồng
func AddContractAccess(c *fiber.Ctx) error {
	// Dữ liệu đầu vào từ request body
	var data struct {
		AccountID  string `json:"accountID"`
		ContractID string `json:"contractID"`
	}

	// Parse request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Gọi hàm service để xử lý logic thêm quyền truy cập
	if err := services.AddContractAccessService(data.AccountID, data.ContractID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(), //trả lỗi từ service
		})
	}

	// Trả về kết quả thành công
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm quyền truy cập thành công",
	})
}
