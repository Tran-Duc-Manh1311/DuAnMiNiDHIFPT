package controllers

import (
	"MiniHIFPT/models"
	"MiniHIFPT/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Hàm đăng ký tài khoản
func Register(c *fiber.Ctx) error {
	// Lấy dữ liệu đầu vào từ request
	var newAccount models.Accounts
	if err := c.BodyParser(&newAccount); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Gọi service để xử lý đăng ký
	// result,errService:= service.Register(&newAccount)
	// if errService !=nil{
	// 	return c.JSON(fiber.Map{
	// 			"status":errService.Status,
	// 			"error": errService.Msg,
	// 		})
	// }else{
	// 	return c.JSON(fiber.Map{
	// 		"status":1,
	// 			"mes": "Ok",
	// 			"data":result,
	// 		})
	// }
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
