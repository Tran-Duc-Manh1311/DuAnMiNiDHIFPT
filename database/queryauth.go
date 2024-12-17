package database

import (
	"MiniHIFPT/models"
	"errors"
	"gorm.io/gorm"
	"regexp"
	// "strings"
	"fmt"
	"github.com/mssola/user_agent"
	// "time"
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

///

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
	// Kiểm tra xem UUID đã tồn tại chưa
	var existingDevice models.Devices
	err := DB.Where("id_uuid = ?", device.ID).First(&existingDevice).Error
	if err == nil {
		// Nếu UUID đã tồn tại, trả về lỗi
		return errors.New("device UUID already exists")
	} else if err != gorm.ErrRecordNotFound {
		// Nếu có lỗi khác ngoài lỗi không tìm thấy bản ghi
		return err
	}

	// Tạo thiết bị mới nếu UUID chưa tồn tại
	return DB.Create(device).Error
}

// Cập nhật thông tin thiết bị trong cơ sở dữ liệu
func UpdateDevice(device *models.Devices) error {
	// Kiểm tra xem thiết bị có tồn tại không
	var existingDevice models.Devices
	err := DB.Where("id_uuid = ?", device.ID).First(&existingDevice).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không tìm thấy thiết bị, trả về lỗi
		return errors.New("device not found")
	}
	// Cập nhật thông tin thiết bị
	return DB.Save(device).Error
}

// Hàm lấy thông tin số lần đăng nhập trong ngày của người dùng

func GetDailyLoginAttempts(soDienThoai string) (*models.LoginAttempt, error) {
	var loginAttempt models.LoginAttempt
	err := DB.Where("SoDienThoai = ? AND DATE(Ngay) = CURRENT_DATE", soDienThoai).First(&loginAttempt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &loginAttempt, err
}

// Hàm lưu số lần đăng nhập thất bại của người dùng
func SaveLoginAttempt(attempt *models.LoginAttempt) error {
	if err := DB.Save(attempt).Error; err != nil {
		return fmt.Errorf("lỗi khi lưu đăng nhập thất bại: %v", err)
	}
	return nil
}

// Lưu thông tin người dùng và thiết bị
func SaveDeviceUser(account *models.Accounts, device *models.Devices) error {
	// Lưu thông tin người dùng vào cơ sở dữ liệu
	err := DB.Create(account).Error
	if err != nil {
		return err
	}

	// Lưu thông tin thiết bị vào cơ sở dữ liệu
	err = DB.Create(device).Error
	if err != nil {
		return err
	}

	// Trả về nil nếu cả hai thông tin được lưu thành công
	return nil
}
