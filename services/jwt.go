package services

import (
	"MiniHIFPT/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecretKey = []byte("your_secret_key") // Khai báo khóa bí mật cho việc ký và xác thực token

// Hàm tạo JWT cho tài khoản
func GenerateJWT(account *models.Accounts) (string, error) {
	claims := jwt.MapClaims{
		"accountID":   account.ID,
		"phoneNumber": account.SoDienThoai,
		"exp":         time.Now().Add(time.Hour).Unix(), // Token hết hạn sau 1 giờ
	}

	// Tạo token với claims và ký bằng HMAC SHA-256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token bằng jwtSecretKey
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err // Nếu có lỗi khi ký token
	}

	return signedToken, nil // Trả về token đã ký
}
