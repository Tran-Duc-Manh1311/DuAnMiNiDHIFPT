package services

import (
	"MiniHIFPT/models"
	"gorm.io/gorm"
)

type ContractService struct {
	DB *gorm.DB
}

// Hàm thống kê số hợp đồng theo trạng thái
func (service *ContractService) CountContractsByStatus(status string) (int64, error) {
	var count int64
	if err := service.DB.Model(&models.Contract{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
