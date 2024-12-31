package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"math"
	"time"
)

// Hàm xử lý thanh toán
func ProcessPayment(invoiceID string, amount float64, method string, accountID string) error {
	// Lấy thông tin hóa đơn
	invoice, err := database.GetInvoiceByID(invoiceID)
	if err != nil {
		return errors.New("hóa đơn không tồn tại")
	}

	// Kiểm tra trạng thái thanh toán của hóa đơn
	if invoice.PaymentStatus == "Paid" {
		return errors.New("hóa đơn đã được thanh toán")
	}

	// Kiểm tra quyền truy cập của tài khoản đối với hợp đồng
	contractID := invoice.ContractID // Giả sử invoice chứa HopDongID (id của hợp đồng liên quan)
	count, err := database.CheckAccess(accountID, contractID)
	if err != nil {
		return errors.New("lỗi khi kiểm tra quyền truy cập")
	}
	if count == 0 {
		return errors.New("tài khoản không có quyền truy cập vào hợp đồng này")
	}

	// Kiểm tra số tiền thanh toán
	if bills := math.Abs(invoice.Amount - amount); bills > 0 {
		return errors.New("số tiền thanh toán phải bằng số tiền của hóa đơn")
	}

	// Kiểm tra phương thức thanh toán
	isValid, err := database.IsPaymentMethodValid(method)
	if err != nil || !isValid {
		return errors.New("phương thức thanh toán không hợp lệ")
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
		return errors.New("không thể tạo bản ghi thanh toán")
	}

	// Đảm bảo payment.ID có giá trị
	if payment.ID == "" {
		return errors.New("không thể lấy ID của thanh toán")
	}

	// Tạo giao dịch thanh toán
	transaction := models.PaymentTransaction{
		PaymentID:       payment.ID,
		TransactionID:   generateTransactionID(),
		TransactionDate: time.Now(),
		Status:          "Success",
	}
	if err := database.CreatePaymentTransaction(&transaction); err != nil {
		return errors.New("không thể tạo giao dịch thanh toán")
	}

	// Cập nhật trạng thái và ngày thanh toán của hóa đơn
	invoice.PaymentStatus = "Paid"
	now := time.Now()
	invoice.PaidDate = &now // Gán giá trị ngày thanh toán hiện tại (dùng con trỏ)

	if err := database.UpdateInvoice(invoice); err != nil {
		return errors.New("không thể cập nhật hóa đơn sau khi thanh toán")
	}

	return nil
}

// Hàm tạo mã giao dịch duy nhất
func generateTransactionID() string {
	// Tạo một mã giao dịch duy nhất sử dụng timestamp
	// timestamp là YYYYMMDDHHMMSS (năm, tháng, ngày, giờ, phút, giây)
	return "TX-" + time.Now().Format("20060102150405")
}
