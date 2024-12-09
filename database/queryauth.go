package database

import (
	"MiniHIFPT/models"
	"errors"
	"gorm.io/gorm"
	"regexp"
	// "strings"
	"github.com/mssola/user_agent"
	"time"
)

// Kiểm tra nếu số điện thoại đã tồn tại trong hệ thống
func CheckExistingAccount(soDienThoai string) (*models.Accounts, error) {
	var existingAccount models.Accounts
	err := DB.Where("SoDienThoai = ?", soDienThoai).First(&existingAccount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &existingAccount, nil
}

// Tạo tài khoản mới
func CreateAccount(newAccount *models.Accounts) error {
	return DB.Create(newAccount).Error
}

// Lấy tài khoản theo số điện thoại
func GetAccountByPhone(soDienThoai string) (*models.Accounts, error) {
	var account models.Accounts
	result := DB.Where("SoDienThoai = ?", soDienThoai).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

// Lấy thiết bị theo số điện thoại
func GetDeviceByPhone(soDienThoai string) (*models.Devices, error) {
	var device models.Devices
	result := DB.Where("SoDienThoai = ?", soDienThoai).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

// Tạo mã OTP mới
func CreateOTP(otp *models.OTPCode) error {
	return DB.Create(otp).Error
}

// Lấy mã OTP theo số điện thoại và mã OTP
func GetOTPByPhoneAndCode(soDienThoai, otpCode string) (*models.OTPCode, error) {
	var otp models.OTPCode
	result := DB.Where("SoDienThoai = ? AND OTP_Code = ?", soDienThoai, otpCode).First(&otp)
	if result.Error != nil {
		return nil, result.Error
	}
	return &otp, nil
}

// Lưu mã OTP
func SaveOTP(otp *models.OTPCode) error {
	return DB.Save(otp).Error
}

// Lưu thiết bị
func SaveDevice(device *models.Devices) error {
	return DB.Create(device).Error
}

// Lấy thiết bị theo số điện thoại và loại thiết bị
func GetDeviceByPhoneAndType(soDienThoai string, deviceType string) (*models.Devices, error) {
	var device models.Devices
	err := DB.Where("SoDienThoai = ? AND DeviceType = ?", soDienThoai, deviceType).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Không tìm thấy thiết bị
	}
	if err != nil {
		return nil, err // Trả về lỗi nếu có lỗi khác
	}
	return &device, nil // Trả về thiết bị nếu tìm thấy
}

// Hàm phân tích tên thiết bị, chuẩn hóa tên thiết bị
func ParseDeviceName(deviceName string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.ReplaceAllString(deviceName, "")
}

// Hàm phân tích hệ điều hành từ User-Agent
// func ParseOperatingSystem(userAgent string) string {
// 	// Chuyển User-Agent thành chữ thường để xử lý dễ dàng
// 	ua := strings.ToLower(userAgent)

// 	// Kiểm tra các hệ điều hành phổ biến
// 	if matched, _ := regexp.MatchString("windows", ua); matched {
// 		return "Windows"
// 	} else if matched, _ := regexp.MatchString("macintosh|mac os x", ua); matched {
// 		return "macOS"
// 	} else if matched, _ := regexp.MatchString("x11|linux", ua); matched {
// 		return "Linux"
// 	} else if matched, _ := regexp.MatchString("android", ua); matched {
// 		return "Android"
// 	} else if matched, _ := regexp.MatchString("iphone|ipad|ios", ua); matched {
// 		return "iOS"
// 	}

//		return "Unknown OS"
//	}
func ParseOperatingSystem(userAgent string) string {
	ua := user_agent.New(userAgent)
	if ua.OS() != "" {
		platform := ua.OS()
		return platform
	}
	return "Unknown OS"
}

// Tạo thiết bị mới trong cơ sở dữ liệu
func CreateDevice(device *models.Devices) error {
	// Kiểm tra xem thiết bị có tồn tại chưa
	var existingDevice models.Devices
	err := DB.Where("DeviceName = ?", device.DeviceName).First(&existingDevice).Error
	if err == nil {
		// Nếu đã tồn tại thiết bị, trả về lỗi
		return errors.New("device already exists")
	}
	// Tạo thiết bị mới
	return DB.Create(device).Error
}

// Cập nhật thông tin thiết bị trong cơ sở dữ liệu
func UpdateDevice(device *models.Devices) error {
	// Kiểm tra xem thiết bị có tồn tại không
	var existingDevice models.Devices
	err := DB.Where("ID = ?", device.ID).First(&existingDevice).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không tìm thấy thiết bị, trả về lỗi
		return errors.New("device not found")
	}
	// Cập nhật thông tin thiết bị
	return DB.Save(device).Error
}

// Hàm lấy thông tin số lần đăng nhập trong ngày của người dùng
func GetDailyLoginAttempts(phone string) (*models.LoginAttempt, error) {
	var attempt models.LoginAttempt
	err := DB.Where("SoDienThoai = ? AND Ngay >= ?", phone, time.Now().Add(-24*time.Hour)).First(&attempt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &attempt, nil
}

// Hàm lưu số lần đăng nhập thất bại của người dùng
func SaveLoginAttempt(attempt *models.LoginAttempt) error {
	return DB.Save(attempt).Error
}
