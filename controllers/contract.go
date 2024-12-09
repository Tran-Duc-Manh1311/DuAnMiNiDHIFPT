package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"regexp"
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

	// Kiểm tra xem hợp đồng đã tồn tại hay chưa
	existingContract, err := database.FindContractByDetails(contract.TenKhachHang, contract.DiaChi, contract.MaTinh, contract.MaQuanHuyen)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Lỗi khi kiểm tra hợp đồng tồn tại",
		})
	}

	if existingContract != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Hợp đồng đã tồn tại",
		})
	}

	//kiểm tra tên khách hàng
	nameRegex := "^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠƯàáâãèéêìíòóôõùúăđĩũơưạ-ỹẠ-Ỹ ]+$"

	matched, err := regexp.MatchString(nameRegex, contract.TenKhachHang)
	if err != nil || !matched {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tên khách hàng không hợp lệ. Chỉ được nhập chữ tiếng Việt có dấu, cả chữ hoa và chữ thường, cùng khoảng trắng.",
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

func UpdateContract(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	// Kiểm tra quyền truy cập
	var count int64
	idUUID, err := uuid.Parse(contractID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID hợp đồng không hợp lệ: " + err.Error(),
		})
	}

	// Kiểm tra quyền của tài khoản với hợp đồng
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

	// Kiểm tra các trường hợp thiếu thông tin
	if updatedData.TenKhachHang == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tên khách hàng không được để trống",
		})
	}
	if updatedData.DiaChi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Địa chỉ không được để trống",
		})
	}
	if updatedData.MaTinh == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Mã tỉnh không được để trống",
		})
	}
	if updatedData.MaQuanHuyen == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Mã quận huyện không được để trống",
		})
	}

	// Kiểm tra tên khách hàng với regex
	nameRegex := "^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠƯàáâãèéêìíòóôõùúăđĩũơưạ-ỹẠ-Ỹ ]+$"
	matched, err := regexp.MatchString(nameRegex, updatedData.TenKhachHang)
	if err != nil || !matched {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tên khách hàng không hợp lệ. Chỉ được nhập chữ tiếng Việt có dấu, cả chữ hoa và chữ thường, cùng khoảng trắng.",
		})
	}

	// Cập nhật các trường hợp hợp đồng
	updates := map[string]interface{}{
		"TenKhachHang": updatedData.TenKhachHang,
		"DiaChi":       updatedData.DiaChi,
		"MaTinh":       updatedData.MaTinh,
		"MaQuanHuyen":  updatedData.MaQuanHuyen,
		"MaPhuongXa":   updatedData.MaPhuongXa,
		"MaDuong":      updatedData.MaDuong,
		"SoNha":        updatedData.SoNha,
	}

	// Xóa các trường có giá trị rỗng khỏi map updates
	for key, value := range updates {
		if value == "" {
			delete(updates, key)
		}
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
