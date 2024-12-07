package database

import (
	"MiniHIFPT/models"
	"errors"

	"gorm.io/gorm"
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

func GetDailyLoginAttempts(phone string) (*models.LoginAttempt, error) {
	var loginAttempt models.LoginAttempt
	// Truy vấn dữ liệu theo số điện thoại và ngày hiện tại
	err := DB.Where("SoDienThoai = ? AND DATE(Ngay) = CURDATE()", phone).
		Order("id_uuid").First(&loginAttempt).Error

	if err != nil {
		// Kiểm tra lỗi nếu không tìm thấy bản ghi
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Không có bản ghi nào thì trả về nil
		}
		return nil, err // Nếu có lỗi khác, trả về lỗi
	}

	return &loginAttempt, nil
}

// Tạo hoặc cập nhật lần nhập sai
func SaveLoginAttempt(loginAttempts *models.LoginAttempt) error {
	return DB.Save(loginAttempts).Error //lỗi
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
