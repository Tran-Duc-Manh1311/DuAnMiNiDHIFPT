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

// Hàm tạo liên kết hợp đồng mới trong cơ sở dữ liệu
func CreateCustomerContract(ctm_tract *models.Customer_Contractt) error {
	// Thực hiện truy vấn để tạo hợp đồng mới
	result := DB.Create(ctm_tract)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Kiểm tra hợp đồng đã liên kết với số điện thoại chưa
func CheckExistingContractByPhoneAndContract(SoDienThoai, HopDongID string) (bool, error) {
	var count int64
	// Truy vấn xem hợp đồng đã liên kết với số điện thoại này chưa
	if err := DB.Model(&models.Customer_Contractt{}).
		Where("SoDienThoai = ? AND HopDongID = ?", SoDienThoai, HopDongID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
