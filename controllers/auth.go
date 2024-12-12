package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
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

// Login Controller
func Login(c *fiber.Ctx) error {
	var loginCredentials models.Accounts
	if err := c.BodyParser(&loginCredentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}
	//
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

	// Kiểm tra thông tin thiết bị
	currentDeviceType := c.Get("User-Agent")                            // Lấy thông tin từ User-Agent
	deviceName := database.ParseDeviceName(currentDeviceType)           // Sử dụng hàm chuẩn hóa tên thiết bị
	operatingSystem := database.ParseOperatingSystem(currentDeviceType) // Lấy hệ điều hành từ User-Agent

	device, err := database.GetDeviceByPhoneAndType(account.SoDienThoai, currentDeviceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi khi kiểm tra thiết bị",
		})
	}

	// Nếu thiết bị chưa có trong hệ thống, tạo mới
	if device == nil {
		device = &models.Devices{
			SoDienThoai:     account.SoDienThoai,
			DeviceName:      deviceName,
			DeviceType:      currentDeviceType,
			OperatingSystem: operatingSystem,
			LanDungGanNhat:  time.Now(),
		}
		if err := database.CreateDevice(device); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi khi tạo thông tin thiết bị" + err.Error(),
			})
		}
	} else {
		// Cập nhật thông tin thiết bị
		device.DeviceName = deviceName
		device.OperatingSystem = operatingSystem
		device.LanDungGanNhat = time.Now()
		if err := database.UpdateDevice(device); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi khi cập nhật thông tin thiết bị",
			})
		}
	}

	// Nếu thiết bị chưa được xác minh, yêu cầu OTP
	if device.XacThucOTP == false {
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
	if device.XacThucOTP {
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

// func Login(c *fiber.Ctx) error {
// 	var loginCredentials models.Accounts
// 	if err := c.BodyParser(&loginCredentials); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Dữ liệu đầu vào không hợp lệ",
// 		})
// 	}

// 	if loginCredentials.SoDienThoai == "" || loginCredentials.MatKhau == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Số điện thoại và mật khẩu không được để trống",
// 		})
// 	}

// 	token, err := services.LoginService(loginCredentials, c.Get("User-Agent"))
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Đăng nhập thành công",
// 		"token":   token,
// 	})
// }

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
	// Kiểm tra nếu số điện thoại hoặc mật khẩu trống
	if otpRequest.SoDienThoai == "" || otpRequest.OTPCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số điện thoại và mật khẩu không được để trống",
		})
	}

	// Kiểm tra mã OTP từ cơ sở dữ liệu
	otp, err := database.GetOTPByPhoneAndCode(otpRequest.SoDienThoai, otpRequest.OTPCode)
	if err != nil || otp == nil || time.Now().After(otp.HetHan) || otp.DaXacThuc {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Mã OTP hoặc số điện thoại không hợp lệ.",
		})
	}
	// Kiểm tra số điện thoại
	phoneRegex := "^\\d{10,15}$" // Số điện thoại phải từ 10 đến 15 ký tự và chỉ chứa số
	matched, err := regexp.MatchString(phoneRegex, otpRequest.SoDienThoai)
	if err != nil || !matched {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số điện thoại không hợp lệ. Chỉ được chứa số và phải từ 10 đến 15 ký tự.",
		})
	}
	// Đánh dấu mã OTP đã được xác thực
	otp.DaXacThuc = true
	if err := database.SaveOTP(otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi lưu trạng thái OTP",
		})
	}

	account, err := database.GetAccountByPhone(otpRequest.SoDienThoai)
	if err != nil || account == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Tài khoản không tồn tại.",
		})
	}

	// Kiểm tra xem thiết bị đã tồn tại hay chưa
	deviceType := c.Get("User-Agent") // Lấy thông tin thiết bị từ User-Agent
	device, err := database.GetDeviceByPhoneAndType(otpRequest.SoDienThoai, deviceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Có lỗi xảy ra khi kiểm tra thông tin thiết bị",
		})
	}

	if device == nil {
		// Nếu thiết bị chưa tồn tại, tạo bản ghi mới
		device = &models.Devices{
			SoDienThoai: otpRequest.SoDienThoai,
			DeviceType:  deviceType,
			XacThucOTP:  true, // Đánh dấu thiết bị đã xác thực OTP
		}
		if err := database.CreateDevice(device); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi xảy ra khi tạo thiết bị",
			})
		}
	} else {
		// Nếu thiết bị đã tồn tại, cập nhật cột XacThucOTP
		device.XacThucOTP = true
		if err := database.UpdateDevice(device); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Có lỗi xảy ra khi cập nhật thiết bị",
			})
		}
	}

	// Nếu OTP và thiết bị đã xác thực thành công, tạo token
	if device.XacThucOTP {
		token, err := generateJWT(account)
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

// Hàm tạo mã OTP ngẫu nhiên
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}
