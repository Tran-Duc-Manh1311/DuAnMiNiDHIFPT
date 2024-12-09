package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Thêm quyền truy cập cho tài khoản đối với hợp đồng
func AddContractAccess(c *fiber.Ctx) error {
	var data struct {
		AccountID  string `json:"accountID"`
		ContractID string `json:"contractID"`
	}

	// Xử lý dữ liệu đầu vào
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Thêm quyền truy cập vào bảng trung gian
	accountContract := models.Account_Contract{
		AccountID:  data.AccountID,
		ContractID: data.ContractID,
	}

	// Lưu quyền truy cập vào cơ sở dữ liệu
	if err := database.DB.Create(&accountContract).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể thêm quyền truy cập",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm quyền truy cập thành công",
	})
}
