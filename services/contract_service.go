package services

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"regexp"
)

// ServiceResponse đại diện cho cấu trúc của một phản hồi từ dịch vụ (service)
type ServiceResponse struct {
	Code    int
	Message string
	Data    interface{}
}

// hàm này để tạo một đối tượng ServiceResponse dễ dàng và nhanh chóng. Nó trả về con trỏ đến đối tượng ServiceResponse.
func respond(code int, message string, data interface{}) *ServiceResponse {
	return &ServiceResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func GetContracts() *ServiceResponse {
	contracts, err := database.GetContracts()
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Không thể lấy thông tin hợp đồng", nil)
	}
	return respond(fiber.StatusOK, "Lấy danh sách hợp đồng thành công", contracts)
}

func GetContractByID(contractID string, accountID string) *ServiceResponse {
	idUUID, err := uuid.Parse(contractID)
	count, err := database.CheckAccess(accountID, contractID)
	if err != nil {
		return respond(fiber.StatusBadRequest, "ID hợp đồng không hợp lệ", nil)
	}
	//kiểm tra quyền truy cập
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra quyền truy cập", nil)
	}
	if count == 0 {
		return respond(fiber.StatusForbidden, "Bạn không có quyền cập nhật hợp đồng này", nil)
	}

	var contract models.Contract
	if err := database.DB.First(&contract, "id_uuid = ?", idUUID).Error; err != nil {
		return respond(fiber.StatusNotFound, "Hợp đồng không tồn tại", nil)
	}

	return respond(fiber.StatusOK, "Lấy thông tin hợp đồng thành công", &contract)
}

func CreateContractService(c *fiber.Ctx) *ServiceResponse {
	var contract models.Contract
	if err := c.BodyParser(&contract); err != nil {
		return respond(fiber.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", nil)
	}

	if contract.TenKhachHang == "" || contract.DiaChi == "" || contract.MaTinh == "" || contract.MaQuanHuyen == "" {
		return respond(fiber.StatusBadRequest, "Thiếu thông tin cần thiết", nil)
	}

	existingContract, err := database.FindContractByDetails(contract.TenKhachHang)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra hợp đồng", nil)
	}
	if existingContract != nil {
		return respond(fiber.StatusConflict, "Hợp đồng đã tồn tại", nil)
	}

	nameRegex := "^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠƯàáâãèéêìíòóôõùúăđĩũơưạ-ỹẠ-Ỹ ]+$"
	matched, _ := regexp.MatchString(nameRegex, contract.TenKhachHang)
	if !matched {
		return respond(fiber.StatusBadRequest, "Tên khách hàng không hợp lệ", nil)
	}

	if err := database.CreateContract(&contract); err != nil {
		return respond(fiber.StatusInternalServerError, "Không thể lưu hợp đồng", nil)
	}

	return respond(fiber.StatusCreated, "Tạo hợp đồng thành công", &contract)
}

func CheckContractStatusService(c *fiber.Ctx) *ServiceResponse {
	// Lấy contractID từ URL params
	contractID := c.Params("id")
	if contractID == "" {
		return respond(fiber.StatusBadRequest, "ID hợp đồng không hợp lệ", nil)
	}

	// Lấy accountID từ context (do middleware đặt vào)
	accountID, ok := c.Locals("accountID").(string)
	if !ok || accountID == "" {
		return respond(fiber.StatusUnauthorized, "Không thể xác thực người dùng", nil)
	}

	// Kiểm tra quyền truy cập CheckAccess
	accessCount, err := database.CheckAccess(accountID, contractID)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra quyền truy cập", nil)
	}
	if accessCount == 0 {
		return respond(fiber.StatusForbidden, "Bạn không có quyền truy cập hợp đồng này", nil)
	}

	// Truy vấn thông tin hợp đồng từ cơ sở dữ liệu
	contract, err := database.GetContractByID(contractID)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi truy vấn cơ sở dữ liệu", nil)
	}

	// Trả về thông tin hợp đồng
	return respond(fiber.StatusOK, "Trạng thái hợp đồng", map[string]interface{}{
		"Hợp đồng": contract.ID,
		"Status":   contract.Status,
	})
}

func UpdateContract(accountID, contractID string, c *fiber.Ctx) *ServiceResponse {
	count, err := database.CheckAccess(accountID, contractID)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra quyền truy cập", nil)
	}
	if count == 0 {
		return respond(fiber.StatusForbidden, "Bạn không có quyền cập nhật hợp đồng này", nil)
	}

	contract, err := database.GetContractByID(contractID)
	if err != nil {
		return respond(fiber.StatusNotFound, "Hợp đồng không tồn tại", nil)
	}

	var updatedData models.Contract
	if err := c.BodyParser(&updatedData); err != nil {
		return respond(fiber.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", nil)
	}

	if updatedData.TenKhachHang == "" {
		return respond(fiber.StatusBadRequest, "Tên khách hàng không được để trống", nil)
	}

	nameRegex := "^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠƯàáâãèéêìíòóôõùúăđĩũơưạ-ỹẠ-Ỹ ]+$"
	matched, _ := regexp.MatchString(nameRegex, updatedData.TenKhachHang)
	if !matched {
		return respond(fiber.StatusBadRequest, "Tên khách hàng không hợp lệ", nil)
	}

	updates := map[string]interface{}{
		"TenKhachHang": updatedData.TenKhachHang,
		"DiaChi":       updatedData.DiaChi,
		"MaTinh":       updatedData.MaTinh,
		"MaQuanHuyen":  updatedData.MaQuanHuyen,
		"MaPhuongXa":   updatedData.MaPhuongXa,
		"MaDuong":      updatedData.MaDuong,
		"SoNha":        updatedData.SoNha,
	}

	for key, value := range updates {
		if value == "" {
			delete(updates, key)
		}
	}

	if len(updates) > 0 {
		if err := database.UpdateContract(&contract, updates); err != nil {
			return respond(fiber.StatusInternalServerError, "Không thể cập nhật hợp đồng", nil)
		}
	}

	return respond(fiber.StatusOK, "Cập nhật hợp đồng thành công", &contract)
}

func DeleteContract(accountID, contractID string) *ServiceResponse {
	count, err := database.CheckAccess(accountID, contractID)
	if err != nil {
		return respond(fiber.StatusInternalServerError, "Lỗi khi kiểm tra quyền truy cập", nil)
	}
	if count == 0 {
		return respond(fiber.StatusForbidden, "Bạn không có quyền xóa hợp đồng này", nil)
	}

	idUUID, err := uuid.Parse(contractID)
	if err != nil {
		return respond(fiber.StatusBadRequest, "ID hợp đồng không hợp lệ", nil)
	}

	if err := database.DeleteContract(idUUID); err != nil {
		return respond(fiber.StatusInternalServerError, "Không thể xóa hợp đồng", nil)
	}

	return respond(fiber.StatusOK, "Xóa hợp đồng thành công", nil)
}
