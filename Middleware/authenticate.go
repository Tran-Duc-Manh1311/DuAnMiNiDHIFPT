package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

var jwtSecretKey = []byte("your_secret_key") // Khóa bí mật JWT

// Authenticate là middleware để xác thực JWT và lưu accountID vào context
func Authenticate(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token không hợp lệ",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token không hợp lệ",
		})
	}

	// Lưu thông tin tài khoản vào context
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Xác nhận quyền sở hữu token không hợp lệ",
		})
	}
	//
	accountID, ok := claims["accountID"].(string)
	if !ok || accountID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "ID tài khoản không hợp lệ hoặc thiếu trong mã thông báo",
		})
	}

	c.Locals("accountID", accountID)
	return c.Next()
}
