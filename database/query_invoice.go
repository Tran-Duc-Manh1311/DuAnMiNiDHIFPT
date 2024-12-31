package database

import (
	"MiniHIFPT/models"
)

// lấy hóa đơn theo ID
func GetInvoiceByID(invoiceID string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := DB.Where("id = ?", invoiceID).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
