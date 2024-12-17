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

	// Kiểm tra nếu số điện thoại đã liên kết với hợp đồng này chưa
	exists, err := database.CheckExistingContractByPhoneAndContract(ctmContract.SoDienThoai, ctmContract.HopDongID)
	if err != nil {
		return errors.New("Lỗi khi kiểm tra hợp đồng liên kết")
	}

	if exists {
		return errors.New("Số điện thoại đã được liên kết với hợp đồng này")
	}

	// Thêm vào bảng customer_contract
	if err := database.CreateCustomerContract(&ctmContract); err != nil {
		return err
	}

	return nil
}

//Mã lỗi 1452 vi phạm ràng buộc khóa ngoại (foreign key constraint).
