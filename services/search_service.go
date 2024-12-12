package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"fmt"
)

// tìm kiếm khách hàng và hợp đồng
func SearchContractsAndCustomers(searchTerm string) ([]models.Contract, []models.Customer, error) {
	// Tìm kiếm hợp đồng
	contracts, err := database.SearchContracts(searchTerm)
	if err != nil {
		return nil, nil, fmt.Errorf("Lỗi khi tìm kiếm hợp đồng: %w", err)
	}

	// Tìm kiếm khách hàng
	customers, err := database.SearchCustomers(searchTerm)
	if err != nil {
		return nil, nil, fmt.Errorf("Lỗi khi tìm kiếm khách hàng: %w", err)
	}

	return contracts, customers, nil
}
