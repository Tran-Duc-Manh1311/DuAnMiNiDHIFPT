package database

import (
	"MiniHIFPT/models"
)

// CountContractsByStatus đếm số hợp đồng theo trạng thái cụ thể
func CountContractsByStatus(status string) (int64, error) {
	var count int64
	// Thực hiện truy vấn GORM
	if err := DB.Model(&models.Contract{}).
		Where("Status = ?", status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
