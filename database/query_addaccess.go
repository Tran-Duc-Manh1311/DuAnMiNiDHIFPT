package database

import (
	"MiniHIFPT/models"
	"errors"
	"gorm.io/gorm"
)

// CreateContractAccess thêm quyền truy cập vào bảng trung gian
func CreateContractAccess(accountContract *models.Account_Contract) error {
	var existing models.Account_Contract

	// Kiểm tra xem accountID và contractID đã tồn tại chưa
	if err := DB.Where("AccountID = ? AND ContractID = ?", accountContract.AccountID, accountContract.ContractID).
		First(&existing).Error; err != nil {

		// Nếu lỗi không phải record not found trả về lỗi
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return errors.New("Quyền truy cập đã tồn tại")
	}

	// Nếu chưa tồn tại, thêm mới vào cơ sở dữ liệu
	if err := DB.Create(accountContract).Error; err != nil {
		return err
	}

	return nil
}
