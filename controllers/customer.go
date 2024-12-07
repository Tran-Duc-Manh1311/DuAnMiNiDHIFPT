package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Lấy thông tin các kh
func GetCustomers(c *fiber.Ctx) error {
	customers, err := database.GetCustomers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin khách hàng",
		})
	}
	return c.JSON(customers)
}
func CreateCustomers(c *fiber.Ctx) error {
	// Tạo một struct tạm thời để nhận dữ liệu JSON
	type TempCustomer struct {
		SoDienThoai   string `json:"SoDienThoai"`
		TenKhachHang  string `json:"TenKhachHang"`
		GioiTinh      string `json:"GioiTinh"`
		NgaySinh      string `json:"NgaySinh"` // Ngày sinh là chuỗi tạm thời
		Email         string `json:"Email"`
		LoaiKhachHang string `json:"LoaiKhachHang"`
	}

	var tempCustomer TempCustomer

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&tempCustomer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
			// "details": err.Error(),
		})
	}
	//...
	// Chuyển đổi ngày sinh từ chuỗi thành kiểu `time.Time`
	var parsedDate *time.Time
	if tempCustomer.NgaySinh != "" {
		date, err := time.Parse("2006-01-02", tempCustomer.NgaySinh)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Định dạng ngày sinh không hợp lệ, yêu cầu dạng YYYY-MM-DD",
				// "details": err.Error(),
			})
		}
		parsedDate = &date
	}

	// Gán dữ liệu vào struct Customer
	customer := models.Customer{
		SoDienThoai:   tempCustomer.SoDienThoai,
		TenKhachHang:  tempCustomer.TenKhachHang,
		GioiTinh:      tempCustomer.GioiTinh,
		NgaySinh:      parsedDate,
		Email:         tempCustomer.Email,
		LoaiKhachHang: tempCustomer.LoaiKhachHang,
	}

	// Kiểm tra dữ liệu khách hàng
	if customer.SoDienThoai == "" || customer.TenKhachHang == "" || customer.GioiTinh == "" || customer.Email == "" || customer.LoaiKhachHang == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	if err := database.CreateCustomer(&customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo khách hàng",
		})
	}
	// Trả về khách hàng đã được thêm vào
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm khách hàng thành công",
		// "data":    customer,
	})
}
