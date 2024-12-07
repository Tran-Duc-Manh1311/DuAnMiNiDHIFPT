package database

import (
	"MiniHIFPT/models"
)

// Lấy tất cả các hợp đồng
func GetCtm_contract() ([]models.Customer_Contractt, error) {
	var ctm_tracts []models.Customer_Contractt
	result := DB.Find(&ctm_tracts)
	return ctm_tracts, result.Error
}

// Hàm tạo hợp đồng mới trong cơ sở dữ liệu
func CreateCustomerContract(ctm_tract *models.Customer_Contractt) error {
	// Thực hiện truy vấn để tạo hợp đồng mới
	result := DB.Create(ctm_tract)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
