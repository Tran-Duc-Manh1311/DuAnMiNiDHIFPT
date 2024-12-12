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

// Kiểm tra xem số điện thoại đã được liên kết với hợp đồng trước đó chưa
func CheckExistingContract(phone string) (int64, error) {
	var count int64
	if err := DB.Model(&models.Customer_Contractt{}).Where("SoDienThoai = ?", phone).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
