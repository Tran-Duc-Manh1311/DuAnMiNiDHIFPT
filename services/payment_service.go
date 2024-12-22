package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"time"
)

// Hàm xử lý thanh toán
func ProcessPayment(invoiceID string, amount float64, method string) error {
	// Kiểm tra hóa đơn
	invoice, err := database.GetInvoiceByID(invoiceID)
	if err != nil {
		return errors.New("Hóa đơn không tồn tại")
	}

	if invoice.PaymentStatus == "Paid" {
		return errors.New("Hóa đơn đã được thanh toán")
	}

	if invoice.Amount < amount {
		return errors.New("Số tiền thanh toán vượt quá số tiền của hóa đơn")
	}

	// Kiểm tra phương thức thanh toán
	isValid, err := database.IsPaymentMethodValid(method)
	if err != nil || !isValid {
		return errors.New("Phương thức thanh toán không hợp lệ")
	}

	// Tạo bản ghi thanh toán
	payment := models.Payment{
		InvoiceID:   invoiceID,
		Amount:      amount,
		Method:      method,
		PaymentDate: time.Now(),
		Status:      "Completed",
	}
	if err := database.CreatePayment(&payment); err != nil {
		return errors.New("Không thể tạo bản ghi thanh toán")
	}

	// Kiểm tra xem payment.ID đã được gán chưa
	if payment.ID == "" {
		return errors.New("Không thể lấy ID của thanh toán")
	}

	// Cập nhật trạng thái hóa đơn
	remainingAmount := invoice.Amount - amount
	if remainingAmount <= 0 {
		invoice.PaymentStatus = "Paid"
		now := time.Now()
		invoice.PaidDate = &now
	} else {
		invoice.PaymentStatus = "Pending"
	}
	if err := database.UpdateInvoice(invoice); err != nil {
		return errors.New("Không thể cập nhật trạng thái hóa đơn")
	}

	// Tạo giao dịch thanh toán
	transaction := models.PaymentTransaction{
		PaymentID:       payment.ID, // Sử dụng ID của bản ghi thanh toán
		TransactionID:   generateTransactionID(),
		TransactionDate: time.Now(),
		Status:          "Success",
	}
	if err := database.CreatePaymentTransaction(&transaction); err != nil {
		return errors.New("Không thể tạo giao dịch thanh toán")
	}

	return nil
}

func generateTransactionID() string {
	// Tạo một mã giao dịch duy nhất sử dụng timestamp
	return "TX-" + time.Now().Format("20060102150405")
}
