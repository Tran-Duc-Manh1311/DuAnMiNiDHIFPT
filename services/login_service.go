package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"MiniHIFPT/security"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

func HandleLogin(c *fiber.Ctx) error {
	var loginCredentials models.Accounts
	if err := c.BodyParser(&loginCredentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	if loginCredentials.SoDienThoai == "" || loginCredentials.MatKhau == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số điện thoại và mật khẩu không được để trống",
		})
	}
	// Kiểm tra số điện thoại hợp lệ
	phoneRegex := "^\\d{10,15}$" // Số điện thoại từ 10 đến 15 ký tự
	matched, err := regexp.MatchString(phoneRegex, loginCredentials.SoDienThoai)
	if err != nil || !matched {
		return fiber.NewError(fiber.StatusBadRequest, "Số điện thoại không hợp lệ")
	}
	account, err := database.GetAccountByPhone(loginCredentials.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Số điện thoại tài khoản không đúng",
		})
	}

	loginAttempts, err := database.GetDailyLoginAttempts(loginCredentials.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi kiểm tra số lần đăng nhập",
		})
	}

	if loginAttempts != nil {
		if time.Since(loginAttempts.Ngay) >= time.Minute {
			loginAttempts.SoLanSai = 0
			loginAttempts.Ngay = time.Now()
			database.SaveLoginAttempt(loginAttempts)
		}
		if loginAttempts.SoLanSai >= 4 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Bạn đã nhập sai mật khẩu quá số lần quy định. Vui lòng thử lại sau 1 phút.",
			})
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.MatKhau), []byte(loginCredentials.MatKhau)); err != nil {
		if loginAttempts == nil {
			loginAttempts = &models.LoginAttempt{
				SoDienThoai: loginCredentials.SoDienThoai,
				SoLanSai:    1,
				Ngay:        time.Now(),
			}
		} else {
			loginAttempts.SoLanSai++
			loginAttempts.Ngay = time.Now()
		}
		database.SaveLoginAttempt(loginAttempts)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": " Mật khẩu tài khoản không đúng",
		})
	}

	if loginAttempts != nil {
		loginAttempts.SoLanSai = 0
		loginAttempts.Ngay = time.Now()
		database.SaveLoginAttempt(loginAttempts)
	}

	currentDeviceType := c.Get("User-Agent")
	deviceName := database.ParseDeviceName(currentDeviceType)
	operatingSystem := database.ParseOperatingSystem(currentDeviceType)

	device, err := database.GetDeviceByPhoneAndType(account.SoDienThoai, currentDeviceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi kiểm tra thiết bị",
		})
	}

	if device == nil {
		device = &models.Devices{
			SoDienThoai:     account.SoDienThoai,
			DeviceName:      deviceName,
			DeviceType:      currentDeviceType,
			OperatingSystem: operatingSystem,
			LanDungGanNhat:  time.Now(),
		}
		database.CreateDevice(device)
	} else {
		device.DeviceName = deviceName
		device.OperatingSystem = operatingSystem
		device.LanDungGanNhat = time.Now()
		database.UpdateDevice(device)
	}

	if !device.XacThucOTP {
		otpCode := security.GenerateOTP()
		database.CreateOTP(&models.OTPCode{
			SoDienThoai: account.SoDienThoai,
			OTP_Code:    otpCode,
			HetHan:      time.Now().Add(5 * time.Minute),
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Mã OTP đã được gửi. Vui lòng nhập mã OTP.",
			"otp_code": otpCode,
		})
	}

	token, err := security.GenerateJWT(account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi tạo token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Đăng nhập thành công, không cần nhập OTP.",
		"token":   token,
	})
}
