package models

import (
	"time"
)

// ----------------------- Khách hàng -----------------------
type TempCustomer struct {
	SoDienThoai   string `json:"SoDienThoai"`
	TenKhachHang  string `json:"TenKhachHang"`
	GioiTinh      string `json:"GioiTinh"`
	NgaySinh      string `json:"NgaySinh"`
	Email         string `json:"Email"`
	LoaiKhachHang string `json:"LoaiKhachHang"`
}
type Customer struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai   string     `json:"soDienThoai" gorm:"column:SoDienThoai"`                          // Số điện thoại khách hàng
	TenKhachHang  string     `json:"tenKhachHang" gorm:"column:TenKhachHang"`                        // Tên khách hàng
	GioiTinh      string     `json:"gioiTinh" gorm:"column:GioiTinh"`                                // Giới tính khách hàng
	NgaySinh      *time.Time `gorm:"type:date;column:NgaySinh" json:"ngaySinh"`
	Email         string     `json:"email" gorm:"column:Email"`                                          // Email khách hàng
	LoaiKhachHang string     `gorm:"type:char(1);default:'T';column:LoaiKhachHang" json:"loaiKhachHang"` // Loại khách hàng: Tiềm năng (T) hoặc Sử dụng dịch vụ (S)
}

// Chỉ định tên bảng
func (Customer) TableName() string {
	return "Customer" // Tên bảng thực tế trong cơ sở dữ liệu
}

// ----------------------- Hợp đồng -----------------------
type Contract struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	TenKhachHang string `gorm:"column:TenKhachHang;not null"`
	DiaChi       string `gorm:"column:DiaChi;not null"`
	MaTinh       string `gorm:"column:MaTinh;not null"`
	MaQuanHuyen  string `gorm:"column:MaQuanHuyen;not null"`
	MaPhuongXa   string `gorm:"column:MaPhuongXa;not null"`
	MaDuong      string `gorm:"column:MaDuong;null"`
	SoNha        string `gorm:"column:SoNha;null"`
	Status       string `gorm:"column:Status;default:'Active'"`
}

// Chỉ định tên bảng trong MySQL
func (Contract) TableName() string {
	return "contractt" // Tên bảng trong cơ sở dữ liệu
}

// ----------------------- Bảng trung gian -----------------------
type Customer_Contractt struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai string `json:"soDienThoai" gorm:"unique;not null;column:SoDienThoai"`          // Số điện thoại (duy nhất)
	HopDongID   string `gorm:"index;not null;column:HopDongID"`                                // ID hợp đồng

}

// Chỉ định tên bảng trong MySQL
func (Customer_Contractt) TableName() string {
	return "customer_contractt" // Tên bảng trong cơ sở dữ liệu
}

// ----------------------- Bảng trung gian quyền truy cập tài khoản -----------------------
type Account_Contract struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	AccountID  string `gorm:"index;not null;column:AccountID"`                                // ID tài khoản
	ContractID string `gorm:"index;not null;column:ContractID"`                               // ID hợp đồng
}

// Chỉ định tên bảng trong MySQL
func (Account_Contract) TableName() string {
	return "accounts_contracts" // Tên bảng trong cơ sở dữ liệu
}

// ----------------------- Tài khoản người dùng -----------------------
type Accounts struct {
	ID                 string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"`     // UUID tự động
	SoDienThoai        string     `json:"soDienThoai" gorm:"unique;not null;column:SoDienThoai"`              // Số điện thoại (duy nhất)
	MatKhau            string     `json:"matKhau" gorm:"not null;column:MatKhau"`                             // Mật khẩu (bắt buộc)
	NgayTao            time.Time  `json:"ngayTao" gorm:"autoCreateTime;column:NgayTao"`                       // Ngày tạo tài khoản
	NgayCapNhat        time.Time  `json:"ngayCapNhat" gorm:"autoUpdateTime;column:NgayCapNhat"`               // Ngày cập nhật tài khoản
	LanDangNhapGanNhat *time.Time `json:"lanDangNhapGanNhat" gorm:"autoUpdateTime;column:LanDangNhapGanNhat"` // Lần đăng nhập gần nhất
}

// bảng LoginAttempts
type LoginAttempt struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"`
	SoDienThoai string    `gorm:"column:SoDienThoai;not null"`
	SoLanSai    int       `gorm:"column:SoLanSai;default:0"`
	Ngay        time.Time `gorm:"column:Ngay;default:current_timestamp"`
	KhoiPhuc    bool      `gorm:"column:KhoiPhuc;default:false"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}

type OTPCode struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai string    `json:"soDienThoai" gorm:"not null;index;column:SoDienThoai"`           // Số điện thoại (có index)
	OTP_Code    string    `json:"otpCode" gorm:"not null;column:OTP_Code"`                        // Mã OTP (bắt buộc)
	NgayTao     time.Time `json:"ngayTao" gorm:"autoCreateTime;column:NgayTao"`                   // Thời gian tạo OTP
	HetHan      time.Time `json:"hetHan" gorm:"not null;column:HetHan"`                           // Thời gian hết hạn OTP
	DaXacThuc   bool      `json:"daXacThuc" gorm:"default:false;column:DaXacThuc"`                // Trạng thái xác thực OTP (mặc định false)
}

type Devices struct {
	ID              string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai     string    `json:"soDienThoai" gorm:"not null;index;column:SoDienThoai"`           // Số điện thoại
	DeviceName      string    `json:"deviceName" gorm:"column:DeviceName;size:100"`                   // Tên thiết bị
	DeviceType      string    `json:"deviceType" gorm:"column:DeviceType;size:50"`                    // Loại thiết bị
	OperatingSystem string    `json:"operatingSystem" gorm:"column:OperatingSystem;size:50"`          // Hệ điều hành
	Status          string    `json:"status" gorm:"column:Status;size:20;default:'Active'"`           // Trạng thái thiết bị
	LanDungGanNhat  time.Time `json:"lanDungGanNhat" gorm:"autoUpdateTime;column:LanDungGanNhat"`     // Lần sử dụng gần nhất
	CreatedAt       time.Time `gorm:"autoCreateTime;column:CreatedAt"`                                // Thời gian tạo
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:UpdatedAt"`                                // Thời gian cập nhật
	XacThucOTP      bool      `json:"xacThucOTP" gorm:"default:false;column:XacThucOTP"`
}

type Invoice struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`                     // UUID tự động
	ContractID    string     `gorm:"not null;column:contract_id"`                                                   // Liên kết với Contract
	Amount        float64    `gorm:"not null;column:amount"`                                                        // Số tiền hóa đơn
	PaymentStatus string     `gorm:"type:enum('Pending','Paid','Overdue');default:'Pending';column:payment_status"` // Trạng thái thanh toán
	DueDate       time.Time  `gorm:"not null;column:due_date"`                                                      // Ngày đến hạn
	PaidDate      *time.Time `gorm:"column:paid_date"`                                                              // Ngày thanh toán (có thể null)
	CreatedAt     time.Time  `gorm:"autoCreateTime;column:created_at"`                                              // Thời gian tạo
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;column:updated_at"`                                              // Thời gian cập nhật
	ServiceName   string     `json:"servicename" gorm:"column:service_name"`                                        // Tên dịch vụ
}

func (Invoice) TableName() string {
	return "invoice" // Tên bảng trong cơ sở dữ liệu
}

type Payment struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`        // UUID tự động
	InvoiceID   string    `gorm:"not null;column:invoice_id"`                                       // Liên kết với Invoice
	Amount      float64   `gorm:"not null;column:amount"`                                           // Số tiền thanh toán
	PaymentDate time.Time `gorm:"autoCreateTime;column:payment_date"`                               // Ngày thanh toán
	Method      string    `gorm:"not null;column:method"`                                           // Phương thức thanh toán
	Status      string    `gorm:"type:enum('Completed','Pending');default:'Pending';column:status"` // Trạng thái thanh toán
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`                                 // Thời gian tạo
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`                                 // Thời gian cập nhật
}

func (Payment) TableName() string {
	return "payment" // Tên bảng trong cơ sở dữ liệu
}

type PaymentMethod struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"` // UUID tự động
	Method    string    `gorm:"unique;not null;column:method"`                             // Mã phương thức thanh toán (duy nhất)
	Name      string    `gorm:"not null;column:name"`                                      // Tên phương thức thanh toán
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`                          // Thời gian tạo
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`                          // Thời gian cập nhật
}

func (PaymentMethod) TableName() string {
	return "paymentmethod" // Tên bảng trong cơ sở dữ liệu
}

type PaymentTransaction struct {
	ID              string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`     // UUID tự động
	PaymentID       string    `gorm:"not null;column:payment_id"`                                    // Liên kết với Payment
	TransactionID   string    `gorm:"unique;not null;column:transaction_id"`                         // Mã giao dịch
	TransactionDate time.Time `gorm:"autoCreateTime;column:transaction_date"`                        // Ngày thực hiện giao dịch
	Status          string    `gorm:"type:enum('Success','Failed');default:'Success';column:status"` // Trạng thái giao dịch
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at"`                              // Thời gian tạo
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at"`                              // Thời gian cập nhật
}

func (PaymentTransaction) TableName() string {
	return "paymenttransaction" // Tên bảng trong cơ sở dữ liệu
}
