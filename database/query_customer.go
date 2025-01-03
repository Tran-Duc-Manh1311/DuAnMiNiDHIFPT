package database

import (
	"MiniHIFPT/models"
)

// Lấy tất cả các hợp đồng
func GetCustomers() ([]models.Customer, error) {
	var customer []models.Customer
	result := DB.Find(&customer)
	return customer, result.Error
}

func GetCustomerByID(customerID string) (*models.Customer, error) {
	var customer models.Customer
	// Tìm khách hàng theo id_uuid
	err := DB.Where("id_uuid = ?", customerID).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func CreateCustomer(customer *models.Customer) error {
	// Thực hiện truy vấn để tạo hợp đồng mới
	result := DB.Create(customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
