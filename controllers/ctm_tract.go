package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các hợp đồng

func Getctm_tracts(c *fiber.Ctx) error {
	ctm_tracts, err := database.GetCtm_contract()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin ",
		})
	}
	return c.JSON(ctm_tracts)
}

// Hàm tạo liên kết hợp đồng
func Createctm_tracts(c *fiber.Ctx) error {

	// Khai báo cấu trúc Customer_Contractt
	var ctm_tract models.Customer_Contractt

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&ctm_tract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra dữ liệu hợp đồng, ví dụ kiểm tra trường hợp hợp đồng đã có
	if ctm_tract.SoDienThoai == "" || ctm_tract.HopDongID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	// Kiểm tra xem số điện thoại đã được liên kết với hợp đồng trước đó chưa
	var existingContractCount int64
	if err := database.DB.Model(&models.Customer_Contractt{}).Where("SoDienThoai = ?", ctm_tract.SoDienThoai).Count(&existingContractCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Lỗi khi kiểm tra số điện thoại",
		})
	}

	// Nếu số điện thoại đã được liên kết với hợp đồng trước đó, trả về lỗi
	if existingContractCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Số điện thoại đã được liên kết với hợp đồng trước đó",
		})
	}

	// Thêm hợp đồng vào cơ sở dữ liệu thông qua hàm query
	if err := database.CreateCustomerContract(&ctm_tract); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo liên kết hợp đồng",
		})
	}

	// Trả về thông báo thành công
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo liên kết số điện thoại và hợp đồng thành công",
	})
}
