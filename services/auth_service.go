package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	// "time"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrPhoneExists  = errors.New("phone number already exists")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)

func RegisterService(newAccount *models.Accounts) error {
	// Kiểm tra dữ liệu đầu vào
	if newAccount.SoDienThoai == "" || newAccount.MatKhau == "" {
		return ErrInvalidInput
	}

	// Kiểm tra định dạng số điện thoại
	phoneRegex := "^\\d{10,15}$"
	matched, err := regexp.MatchString(phoneRegex, newAccount.SoDienThoai)
	if err != nil || !matched {
		return ErrInvalidInput
	}

	// Kiểm tra nếu tài khoản đã tồn tại
	existingAccount, err := database.CheckExistingAccount(newAccount.SoDienThoai)
	if err != nil {
		return ErrInternal
	}
	if existingAccount != nil {
		return ErrPhoneExists
	}

	// Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAccount.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternal
	}
	newAccount.MatKhau = string(hashedPassword)

	// Lưu vào cơ sở dữ liệu
	if err := database.CreateAccount(newAccount); err != nil {
		return ErrInternal
	}

	return nil
}
