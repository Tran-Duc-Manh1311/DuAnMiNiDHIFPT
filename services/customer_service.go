package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"regexp"
	"time"
)

func GetAllCustomers() ([]models.Customer, error) {
	return database.GetCustomers()
}
func GetAllCustomerByID(customerID string) *ServiceResponse {
	// Chuyển đổi customerID thành UUID
	idUUID, err := uuid.Parse(customerID)
	if err != nil {
		return respond(fiber.StatusBadRequest, "ID khách hàng không hợp lệ", nil)
	}

	// Gọi hàm GetCustomerByID để lấy thông tin khách hàng
	customer, err := database.GetCustomerByID(idUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respond(fiber.StatusNotFound, "Khách hàng không tồn tại", nil)
		}
		return respond(fiber.StatusInternalServerError, "Lỗi khi lấy thông tin khách hàng", nil)
	}

	return respond(fiber.StatusOK, "Lấy thông tin chi tiết khách hàng thành công", customer)
}

func CreateCustomerService(tempCustomer *models.TempCustomer) (*models.Customer, error) {
	// Chuyển đổi ngày sinh
	var parsedDate *time.Time
	if tempCustomer.NgaySinh != "" {
		date, err := time.Parse("2006-01-02", tempCustomer.NgaySinh)
		if err != nil {
			return nil, errors.New("định dạng ngày sinh không hợp lệ, yêu cầu dạng YYYY-MM-DD")
		}
		parsedDate = &date
	}

	// Gán dữ liệu vào struct Customer
	customer := &models.Customer{
		SoDienThoai:   tempCustomer.SoDienThoai,
		TenKhachHang:  tempCustomer.TenKhachHang,
		GioiTinh:      tempCustomer.GioiTinh,
		NgaySinh:      parsedDate,
		Email:         tempCustomer.Email,
		LoaiKhachHang: tempCustomer.LoaiKhachHang,
	}

	// Kiểm tra dữ liệu cần thiết
	if customer.SoDienThoai == "" || customer.TenKhachHang == "" || customer.GioiTinh == "" || customer.Email == "" || customer.LoaiKhachHang == "" {
		return nil, errors.New("thiếu thông tin cần thiết")
	}

	// Kiểm tra tên khách hàng
	nameRegex := "^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠƯàáâãèéêìíòóôõùúăđĩũơưạ-ỹẠ-Ỹ ]+$"
	matched, err := regexp.MatchString(nameRegex, customer.TenKhachHang)
	if err != nil || !matched {
		return nil, errors.New("tên khách hàng không hợp lệ. Chỉ được nhập chữ tiếng việt có dấu, cả chữ hoa và chữ thường, cùng khoảng trắng")
	}

	// Lưu khách hàng vào DB
	if err := database.CreateCustomer(customer); err != nil {
		return nil, errors.New("không thể tạo khách hàng")
	}

	return customer, nil
}
