package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
)

// AddContractAccessService xử lý logic thêm quyền truy cập
func AddContractAccessService(accountID string, contractID string) error {
	// Tạo struct để lưu vào database
	accountContract := models.Account_Contract{
		AccountID:  accountID,
		ContractID: contractID,
	}

	if err := database.CreateContractAccess(&accountContract); err != nil {
		return err // Trả về lỗi gốc từ Database Layer
	}

	return nil
}
