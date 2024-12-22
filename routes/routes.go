package routes

import (
	"MiniHIFPT/controllers"
	"MiniHIFPT/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// Route lấy thông tin khách hàng (cần xác thực)
	app.Get("/customers", middleware.Authenticate, controllers.GetCustomers)
	// Route thêm khách hàng mới (cần xác thực)
	app.Post("/customers", middleware.Authenticate, controllers.CreateCustomers)
	// Route chuyển hợp đồng (cần xác thực)
	app.Post("/transfer", middleware.Authenticate, controllers.TransferOwnership)
	// Route lấy thông tin hợp đồng (cần xác thực)
	app.Get("/contracts", middleware.Authenticate, controllers.GetContracts)
	// Route lấy thông tin hợp đồng theo ID (cần xác thực)
	app.Get("/contracts/:id", middleware.Authenticate, controllers.GetContractByID)
	//tìm kiếm
	app.Get("/search", middleware.Authenticate, controllers.SearchContractOrCustomer)
	//kiểm tra trạng thái của hợp đồng
	app.Get("status/:id", middleware.Authenticate, controllers.CheckContractStatusHandler)
	//thống kê hợp đồng theo trạng thái
	app.Get("contract/status/:status", middleware.Authenticate, controllers.GetContractCountByStatus)

	// Route tạo hợp đồng mới (cần xác thực)
	app.Post("/contracts", middleware.Authenticate, controllers.CreateContractHandler)
	// Route sửa hợp đồng (cần xác thực)
	app.Put("/contracts/:id", middleware.Authenticate, controllers.UpdateContract)
	// Route xóa hợp đồng (cần xác thực)
	app.Delete("/contracts/:id", middleware.Authenticate, controllers.DeleteContract)
	// Liên kết hợp đồng (cần xác thực)
	app.Get("/ctmtract", middleware.Authenticate, controllers.GetCtmContracts)
	app.Post("/Createctmtract", middleware.Authenticate, controllers.CreateCtmContracts)

	// Route đăng ký tài khoản (không cần xác thực)
	app.Post("/register", controllers.Register)
	// Route đăng nhập (không cần xác thực)
	app.Post("/login", controllers.Login)
	// Route xác thực OTP (không cần xác thực)
	app.Post("/otp", controllers.VerifyOTP)

	// Route cấp quyền cho tài khoản đối với hợp đồng (cần xác thực)
	app.Post("/grant-access", middleware.Authenticate, controllers.AddContractAccess)
	//tạo hóa đơn
	app.Post("/invoices", middleware.Authenticate, controllers.CreateInvoice)
	//thanh toán hóa đơn
	app.Post("/payment", middleware.Authenticate, controllers.ProcessPayment)

}
