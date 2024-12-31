package services

import (
	"MiniHIFPT/database"
	"errors"
)

// Chuyển sở hữu hợp đồng từ khách hàng cũ sang khách hàng mới
func TransferOwnership(accountID, oldCustomerID, newCustomerID string) error {
	// Kiểm tra quyền truy cập của tài khoản đối với khách hàng cũ
	accessCount, err := database.CheckAccess(accountID, oldCustomerID)
	if err != nil {
		return errors.New("Lỗi khi kiểm tra quyền truy cập")
	}
	if accessCount == 0 {
		return errors.New("Bạn không có quyền truy cập ")
	}

	// Kiểm tra xem khách hàng cũ có tồn tại không
	oldCustomer, err := database.FindCustomerByID(oldCustomerID)
	if err != nil {
		return errors.New("Không tìm thấy khách hàng cũ")
	}

	// Kiểm tra xem khách hàng mới có tồn tại không
	newCustomer, err := database.FindCustomerByID(newCustomerID)
	if err != nil {
		return errors.New("Không tìm thấy khách hàng mới")
	}

	// Lấy tất cả hợp đồng của khách hàng cũ dựa trên số điện thoại
	customerContracts, err := database.FindCustomerContractsByPhoneNumber(oldCustomer.SoDienThoai)
	if err != nil {
		return errors.New("Không thể tìm thấy các hợp đồng của khách hàng cũ")
	}

	// Nếu không có hợp đồng nào thuộc khách hàng cũ
	if len(customerContracts) == 0 {
		return errors.New("Không có hợp đồng nào thuộc về khách hàng cũ")
	}

	// Cập nhật tất cả các hợp đồng để chuyển sang khách hàng mới
	for _, contract := range customerContracts {
		// Cập nhật cột `SoDienThoai` sang số điện thoại của khách hàng mới
		if err := database.TransferContractOwnership(&contract, newCustomer.SoDienThoai); err != nil {
			return errors.New("Không thể chuyển nhượng hợp đồng")
		}
	}

	return nil
}
