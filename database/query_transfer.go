package database

import (
	"MiniHIFPT/models"
)

// Hàm tìm khách hàng theo ID
func FindCustomerByID(id string) (*models.Customer, error) {
	var customer models.Customer
	if err := DB.First(&customer, "id_uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// Hàm lấy tất cả hợp đồng của khách hàng dựa trên số điện thoại
func FindCustomerContractsByPhoneNumber(phone string) ([]models.Customer_Contractt, error) {
	var contracts []models.Customer_Contractt
	if err := DB.Where("SoDienThoai = ?", phone).Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

// Hàm cập nhật hợp đồng để chuyển số điện thoại sang khách hàng mới
func TransferContractOwnership(contract *models.Customer_Contractt, newPhone string) error {
	return DB.Model(&contract).Update("SoDienThoai", newPhone).Error
}
