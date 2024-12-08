package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

//viết 1 file server để lưu các ...

// Hàm đăng ký tài khoản
func Register(c *fiber.Ctx) error {
	var newAccount models.Accounts
	if err := c.BodyParser(&newAccount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra tài khoản đã tồn tại
	existingAccount, err := database.CheckExistingAccount(newAccount.SoDienThoai)
	if err != nil {
		// Xử lý lỗi từ cơ sở dữ liệu
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi kiểm tra tài khoản",
		})
	}

	if existingAccount != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Số điện thoại đã tồn tại trong hệ thống",
		})
	}

	// Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAccount.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể mã hóa mật khẩu",
		})
	}
	newAccount.MatKhau = string(hashedPassword)

	// Lưu tài khoản vào cơ sở dữ liệu
	if err := database.CreateAccount(&newAccount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo tài khoản mới",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo tài khoản thành công. Vui lòng đăng nhập.",
	})
}

func Login(c *fiber.Ctx) error {
	var loginCredentials models.Accounts
	if err := c.BodyParser(&loginCredentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra nếu số điện thoại hoặc mật khẩu trống
	if loginCredentials.SoDienThoai == "" || loginCredentials.MatKhau == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số điện thoại và mật khẩu không được để trống",
		})
	}

	// Lấy thông tin tài khoản từ cơ sở dữ liệu
	account, err := database.GetAccountByPhone(loginCredentials.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Số điện thoại hoặc mật khẩu không đúng",
		})
	}

	// Kiểm tra số lần nhập sai
	loginAttempts, err := database.GetDailyLoginAttempts(loginCredentials.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi kiểm tra số lần đăng nhập",
		})
	}

	if loginAttempts != nil {
		// Reset số lần nhập sai nếu quá 1 phút kể từ lần nhập cuối cùng
		if time.Since(loginAttempts.Ngay) >= time.Minute {
			loginAttempts.SoLanSai = 0
			loginAttempts.Ngay = time.Now()
			if err := database.SaveLoginAttempt(loginAttempts); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Có lỗi xảy ra khi lưu số lần đăng nhập",
				})
			}
		}

		// Chặn đăng nhập nếu vượt quá số lần sai
		if loginAttempts.SoLanSai >= 4 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Bạn đã nhập sai mật khẩu quá số lần quy định. Vui lòng thử lại sau 1 phút.",
			})
		}
	}

	// So sánh mật khẩu
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

		if err := database.SaveLoginAttempt(loginAttempts); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi xảy ra khi lưu số lần đăng nhập",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Số điện thoại hoặc mật khẩu không đúng",
		})
	}

	// Reset số lần nhập sai nếu đăng nhập thành công
	if loginAttempts != nil {
		loginAttempts.SoLanSai = 0
		loginAttempts.Ngay = time.Now()
		if err := database.SaveLoginAttempt(loginAttempts); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi xảy ra khi lưu số lần đăng nhập",
			})
		}
	}

	// Kiểm tra OTP và trạng thái thiết bị
	currentDeviceType := c.Get("User-Agent")
	device, err := database.GetDeviceByPhoneAndType(account.SoDienThoai, currentDeviceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi kiểm tra thiết bị",
		})
	}

	// Nếu thiết bị chưa được xác minh, yêu cầu OTP
	if device == nil {
		otpCode := generateOTP()
		otp := &models.OTPCode{
			SoDienThoai: account.SoDienThoai,
			OTP_Code:    otpCode,
			HetHan:      time.Now().Add(5 * time.Minute),
		}
		if err := database.CreateOTP(otp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi khi tạo OTP",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Mã OTP đã được gửi. Vui lòng nhập mã OTP.",
			"otp_code": otpCode,
		})
	}

	// Tạo JWT token nếu thiết bị đã được xác minh
	if device != nil && device.XacThucOTP {
		token, err := generateJWT(account)
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

	// Nếu không phải thiết bị đã được xác thực, yêu cầu OTP
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Yêu cầu mã OTP để hoàn tất đăng nhập.",
	})
}

func VerifyOTP(c *fiber.Ctx) error {
	var otpRequest struct {
		SoDienThoai string `json:"SoDienThoai"`
		OTPCode     string `json:"OTPCode"`
	}

	// Phân tích dữ liệu đầu vào từ yêu cầu
	if err := c.BodyParser(&otpRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra mã OTP từ cơ sở dữ liệu
	otp, err := database.GetOTPByPhoneAndCode(otpRequest.SoDienThoai, otpRequest.OTPCode)
	if err != nil || otp == nil || time.Now().After(otp.HetHan) || otp.DaXacThuc {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Mã OTP không hợp lệ hoặc đã hết hạn.",
		})
	}

	// Đánh dấu mã OTP đã được xác thực
	otp.DaXacThuc = true
	if err := database.SaveOTP(otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi lưu trạng thái OTP",
		})
	}

	// Lưu thiết bị và đánh dấu thiết bị là đáng tin cậy
	device := &models.Devices{
		SoDienThoai: otpRequest.SoDienThoai,
		DeviceType:  c.Get("User-Agent"),
		XacThucOTP:  true, // Đánh dấu thiết bị này là đã xác thực OTP
	}
	if err := database.SaveDevice(device); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi lưu thiết bị",
		})
	}

	// Tạo JWT token
	account, err := database.GetAccountByPhone(otpRequest.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi lấy thông tin tài khoản",
		})
	}

	token, err := generateJWT(account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi tạo token",
		})
	}

	// Trả về phản hồi thành công cùng với JWT token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Xác thực OTP thành công. Thiết bị đã được lưu.",
		"token":   token,
	})
}

// Hàm tạo mã OTP ngẫu nhiên
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}
