package controllers

import (
	"MiniHIFPT/services"
	"github.com/gofiber/fiber/v2"
)

// lấy danh sách hợp đồng
func GetContracts(c *fiber.Ctx) error {
	contracts := services.GetContracts()
	return c.JSON(contracts)
}

// lấy hợp đồng theo ID
func GetContractByID(c *fiber.Ctx) error {
	contractID := c.Params("id")

	accountID := c.Locals("accountID").(string)

	contract := services.GetContractByID(contractID, accountID)
	return c.JSON(contract)
}

// Tạo hợp đồng mới (thêm)
func CreateContractHandler(c *fiber.Ctx) error {
	err := services.CreateContractService(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Message,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo hợp đồng thành công",
	})
}

// kiểm tra trạng thái hợp đồng
func CheckContractStatusHandler(c *fiber.Ctx) error {
	response := services.CheckContractStatusService(c)
	return c.JSON(response)
}

func UpdateContract(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	err := services.UpdateContract(accountID, contractID, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Message})
	}

	return c.JSON(fiber.Map{"message": "Cập nhật hợp đồng thành công"})
}

func DeleteContract(c *fiber.Ctx) error {
	contractID := c.Params("id")
	accountID := c.Locals("accountID").(string)

	err := services.DeleteContract(accountID, contractID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Message})
	}

	return c.JSON(fiber.Map{"message": "Xóa hợp đồng thành công"})
}
