package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	// "errors"
	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"

	"github.com/google/uuid"
)

// Lấy thông tin các hợp đồng
func GetContracts(c *fiber.Ctx) error {
	contracts, err := database.GetContracts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin hợp đồng",
		})
	}
	return c.JSON(contracts)
}

// Thêm quyền truy cập cho tài khoản đối với hợp đồng
func AddContractAccess(c *fiber.Ctx) error {
	var data struct {
		AccountID  string `json:"accountID"`
		ContractID string `json:"contractID"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Thêm quyền truy cập vào bảng trung gian
	accountContract := models.Account_Contract{
		AccountID:  data.AccountID,
		ContractID: data.ContractID,
	}

	if err := database.DB.Create(&accountContract).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể thêm quyền truy cập",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm quyền truy cập thành công",
	})
}

// Lấy hợp đồng theo ID (chỉ cho phép xem hợp đồng của tài khoản)
func GetContractByID(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	// Kiểm tra quyền truy cập
	var count int64
	idUUID, err := uuid.Parse(contractID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"error": "ID hợp đồng không hợp lệ" + err.Error(),
		})
	}

	if err := database.DB.Model(&models.Account_Contract{}).
		Where("AccountID = ? AND ContractID = ?", accountID, idUUID).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Lỗi khi kiểm tra quyền truy cập",
		})
	}

	if count == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bạn không có quyền truy cập hợp đồng này",
		})
	}

	// Lấy thông tin hợp đồng
	var contract models.Contract
	if err := database.DB.First(&contract, "id_uuid = ?", idUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	return c.JSON(contract)
}

// Tạo hợp đồng mới (thêm)
func CreateContract(c *fiber.Ctx) error {
	var contract models.Contract
	if err := c.BodyParser(&contract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra dữ liệu hợp đồng
	if contract.TenKhachHang == "" || contract.DiaChi == "" || contract.MaTinh == "" || contract.MaQuanHuyen == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	if err := database.CreateContract(&contract); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo hợp đồng",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo hợp đồng thành công",
	})
}

// Sửa thông tin hợp đồng
func UpdateContract(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	// Kiểm tra quyền truy cập
	var count int64
	idUUID, err := uuid.Parse(contractID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"error": "ID hợp đồng không hợp lệ" + err.Error(),
		})
	}

	if err := database.DB.Model(&models.Account_Contract{}).
		Where("AccountID = ? AND ContractID = ?", accountID, idUUID).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Lỗi khi kiểm tra quyền truy cập",
		})
	}

	if count == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bạn không có quyền sửa hợp đồng này",
		})
	}
	// Lấy thông tin hợp đồng từ cơ sở dữ liệu
	contract, err := database.GetContractByID(contractID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	// Phân tích dữ liệu đầu vào từ client
	var updatedData models.Contract
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Cập nhật các trường hợp hợp đồng
	updates := map[string]interface{}{}
	if updatedData.TenKhachHang != "" {
		updates["TenKhachHang"] = updatedData.TenKhachHang
	}
	if updatedData.DiaChi != "" {
		updates["DiaChi"] = updatedData.DiaChi
	}
	if updatedData.MaTinh != "" {
		updates["MaTinh"] = updatedData.MaTinh
	}
	if updatedData.MaQuanHuyen != "" {
		updates["MaQuanHuyen"] = updatedData.MaQuanHuyen
	}
	if updatedData.MaPhuongXa != "" {
		updates["MaPhuongXa"] = updatedData.MaPhuongXa
	}
	if updatedData.MaDuong != "" {
		updates["MaDuong"] = updatedData.MaDuong
	}
	if updatedData.SoNha != "" {
		updates["SoNha"] = updatedData.SoNha
	}

	// Nếu có trường hợp cần cập nhật, thực hiện cập nhật
	if len(updates) > 0 {
		if err := database.UpdateContract(&contract, updates); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể cập nhật hợp đồng",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Sửa hợp đồng thành công",
	})
}

// Xóa hợp đồng
func DeleteContract(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	// Kiểm tra quyền truy cập
	var count int64
	idUUID, err := uuid.Parse(contractID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID hợp đồng không hợp lệ",
		})
	}

	if err := database.DB.Model(&models.Account_Contract{}).
		Where("AccountID = ? AND ContractID = ?", accountID, idUUID).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Lỗi khi kiểm tra quyền truy cập",
		})
	}

	if count == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bạn không có quyền xóa hợp đồng này",
		})
	}

	// Xóa hợp đồng
	if err := database.DB.Delete(&models.Contract{}, "id_uuid = ?", idUUID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể xóa hợp đồng",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Xóa hợp đồng thành công",
	})
}
