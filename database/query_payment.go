package database

import (
	"MiniHIFPT/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// tạo phương thức thanh toán(tiền mặt , thẻ)
func CreatePMMethod(method *models.PaymentMethod) error {
	result := DB.Create(&method)
	return result.Error
}

// kiểm tra tính hợp lệ của phương thức thanh toán
func IsPaymentMethodValid(method string) (bool, error) {
	var count int64
	err := DB.Model(&models.PaymentMethod{}).Where("method = ?", method).Count(&count).Error
	return count > 0, err
}

// Kiểm tra hợp đồng có tồn tại dựa trên các thông tin chi tiết
func Methoddetails(namemethod string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := DB.Where("name = ?",
		namemethod).First(&method).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Không tìm thấy phương thức thanh toán
		}
		return nil, err // Lỗi khi truy vấn
	}
	return &method, nil // phương thức thanh toán đã tồn tại
}

// tạo bản ghi thanh toán mới
func CreatePayment(payment *models.Payment) error {
	payment.ID = uuid.New().String() // Tạo UUID
	// Sử dụng DB.Create để lưu bản ghi thanh toán mới vào cơ sở dữ liệu
	if err := DB.Create(payment).Error; err != nil {
		return err
	}
	// Sau khi tạo thành công, ID của payment sẽ được tự động gán bởi GORM
	return nil
}

// tạo giao dịch thanh toán mới
func CreatePaymentTransaction(transaction *models.PaymentTransaction) error {
	return DB.Create(transaction).Error
}

// cập nhật hóa đơn
func UpdateInvoice(invoice *models.Invoice) error {
	return DB.Save(invoice).Error
}
