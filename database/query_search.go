package database

import (
	"MiniHIFPT/models"
)

// Tìm kiếm khách hàng theo số điện thoại hoặc tên khách hàng
func SearchCustomers(searchTerm string) ([]models.Customer, error) {
	var customers []models.Customer
	if err := DB.Where("SoDienThoai LIKE ? OR TenKhachHang LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// Tìm kiếm hợp đồng theo id_uuid
func SearchContracts(searchTerm string) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := DB.Where("id_uuid LIKE ?", "%"+searchTerm+"%").Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}
