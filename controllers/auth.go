package controllers

import (
	// "MiniHIFPT/database"
	"MiniHIFPT/models"
	// "MiniHIFPT/security"
	"MiniHIFPT/services"
	// "math/rand"
	// "regexp"
	// "strconv"
	// "time"
	"github.com/gofiber/fiber/v2"
)

// Hàm đăng ký tài khoản
func Register(c *fiber.Ctx) error {
	// Lấy dữ liệu đầu vào từ request
	var newAccount models.Accounts
	if err := c.BodyParser(&newAccount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Gọi service để xử lý đăng ký
	if err := services.RegisterService(&newAccount); err != nil {
		// Trả về lỗi phù hợp dựa trên loại lỗi từ service
		switch err {
		case services.ErrInvalidInput:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Số điện thoại hoặc mật khẩu không hợp lệ",
			})
		case services.ErrPhoneExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Số điện thoại đã tồn tại",
			})
		case services.ErrInternal:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Lỗi nội bộ khi xử lý yêu cầu",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Lỗi không xác định",
			})
		}
	}

	// Trả về thành công
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo tài khoản thành công. Vui lòng đăng nhập.",
	})
}

func Login(c *fiber.Ctx) error {
	return services.HandleLogin(c)
}

func VerifyOTP(c *fiber.Ctx) error {
	return services.HandleVerifyOTP(c)
}
