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
func CreateCtmContract(ctmContract models.Customer_Contractt, accountID string) error {

	// Kiểm tra dữ liệu hợp đồng
	if ctmContract.SoDienThoai == "" || ctmContract.HopDongID == "" {
		return errors.New("thiếu thông tin cần thiết")
	}

	// Kiểm tra quyền truy cập của tài khoản với hợp đồng
	count, err := database.CheckAccess(accountID, ctmContract.HopDongID)
	if err != nil {
		return errors.New("lỗi khi kiểm tra quyền truy cập hợp đồng")
	}

	// Nếu không có quyền truy cập, trả lỗi
	if count == 0 {
		return errors.New("tài khoản không có quyền truy cập hợp đồng này")
	}

	// Kiểm tra nếu số điện thoại đã liên kết với hợp đồng này chưa
	exists, err := database.CheckExistingContractByPhoneAndContract(ctmContract.SoDienThoai, ctmContract.HopDongID)
	if err != nil {
		return errors.New("lỗi khi kiểm tra hợp đồng liên kết")
	}

	if exists {
		return errors.New("số điện thoại đã được liên kết với hợp đồng này")
	}

	// Thêm vào bảng customer_contract
	if err := database.CreateCustomerContract(&ctmContract); err != nil {
		return err
	}

	return nil
}

//Mã lỗi 1452 vi phạm ràng buộc khóa ngoại (foreign key constraint).
