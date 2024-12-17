package services

import (
	"MiniHIFPT/database"
)

// CountContractsByStatus trả về số lượng hợp đồng dựa trên trạng thái cụ thể
func CountContractsByStatus(status string) (int64, error) {
	// Gọi hàm từ tầng database
	count, err := database.CountContractsByStatus(status)
	if err != nil {
		return 0, err
	}
	return count, nil
}
