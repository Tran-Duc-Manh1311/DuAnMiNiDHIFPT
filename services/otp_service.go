package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"MiniHIFPT/security"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"time"
)

func HandleVerifyOTP(c *fiber.Ctx) error {
	var otpRequest struct {
		SoDienThoai string `json:"SoDienThoai"`
		OTPCode     string `json:"OTPCode"`
	}

	// Phân tích dữ liệu đầu vào từ yêu cầu
	if err := c.BodyParser(&otpRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
	}

	// Kiểm tra dữ liệu rỗng
	if otpRequest.SoDienThoai == "" || otpRequest.OTPCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Số điện thoại và OTP không được để trống")
	}

	// Kiểm tra số điện thoại hợp lệ
	phoneRegex := "^\\d{10,15}$" // Số điện thoại từ 10 đến 15 ký tự
	matched, err := regexp.MatchString(phoneRegex, otpRequest.SoDienThoai)
	if err != nil || !matched {
		return fiber.NewError(fiber.StatusBadRequest, "Số điện thoại không hợp lệ")
	}

	// Kiểm tra mã OTP từ cơ sở dữ liệu
	otp, err := database.GetOTPByPhoneAndCode(otpRequest.SoDienThoai, otpRequest.OTPCode)
	if err != nil || otp == nil || time.Now().After(otp.HetHan) || otp.DaXacThuc {
		return fiber.NewError(fiber.StatusUnauthorized, "Mã OTP hoặc số điện thoại không hợp lệ")
	}

	// Đánh dấu mã OTP đã được xác thực
	otp.DaXacThuc = true
	if err := database.SaveOTP(otp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Có lỗi xảy ra khi lưu trạng thái OTP")
	}

	// Lấy thông tin tài khoản từ số điện thoại
	account, err := database.GetAccountByPhone(otpRequest.SoDienThoai)
	if err != nil || account == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Tài khoản không tồn tại")
	}

	// Kiểm tra thiết bị
	deviceType := c.Get("User-Agent")
	device, err := database.GetDeviceByPhoneAndType(otpRequest.SoDienThoai, deviceType)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Có lỗi xảy ra khi kiểm tra thông tin thiết bị")
	}

	// Nếu thiết bị chưa tồn tại, tạo mới
	if device == nil {
		device = &models.Devices{
			SoDienThoai: otpRequest.SoDienThoai,
			DeviceType:  deviceType,
			XacThucOTP:  true,
		}
		if err := database.CreateDevice(device); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Có lỗi xảy ra khi tạo thiết bị")
		}
	} else {
		// Cập nhật thiết bị nếu đã tồn tại
		device.XacThucOTP = true
		if err := database.UpdateDevice(device); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Có lỗi xảy ra khi cập nhật thiết bị")
		}
	}

	// Nếu OTP và thiết bị đã xác thực thành công, tạo token
	if device.XacThucOTP {
		token, err := security.GenerateJWT(account) // Nhận cả token và lỗi
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi xảy ra khi tạo token",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Xác minh OTP thành công.",
			"token":   token,
		})
	}
	// Trả về lỗi nếu không xác thực OTP và thiết bị
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Không xác thực được thiết bị.",
	})

}
