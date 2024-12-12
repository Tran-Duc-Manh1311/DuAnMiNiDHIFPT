package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
)

// Lấy tất cả hợp đồng
func GetCtmContracts() ([]models.Customer_Contractt, error) {
	ctmContracts, err := database.GetCtm_contract()
	if err != nil {
		return nil, err
	}
	return ctmContracts, nil
}

// Tạo liên kết hợp đồng
func CreateCtmContract(ctmContract models.Customer_Contractt) error {
	// Kiểm tra dữ liệu hợp đồng
	if ctmContract.SoDienThoai == "" || ctmContract.HopDongID == "" {
		return errors.New("Thiếu thông tin cần thiết")
	}

	// Kiểm tra số điện thoại đã được liên kết với hợp đồng trước đó
	existingContractCount, err := database.CheckExistingContract(ctmContract.SoDienThoai)
	if err != nil {
		return err
	}

	// Nếu số điện thoại đã được liên kết với hợp đồng trước đó
	if existingContractCount > 0 {
		return errors.New("Số điện thoại đã được liên kết với hợp đồng trước đó")
	}

	// Thêm hợp đồng vào cơ sở dữ liệu
	if err := database.CreateCustomerContract(&ctmContract); err != nil {
		return err
	}

	return nil
}
