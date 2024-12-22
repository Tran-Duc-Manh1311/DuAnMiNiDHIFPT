package database

import (
	"MiniHIFPT/models"
)

// lấy hóa đơn theo ID
func GetInvoiceByID(invoiceID string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := DB.Where("id = ?", invoiceID).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// kiểm tra tính hợp lệ của phương thức thanh toán
func IsPaymentMethodValid(method string) (bool, error) {
	var count int64
	err := DB.Model(&models.PaymentMethod{}).Where("method = ?", method).Count(&count).Error
	return count > 0, err
}

// tạo bản ghi thanh toán mới
func CreatePayment(payment *models.Payment) error {
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
