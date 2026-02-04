package wallet

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	wallet_service "github.com/shuTwT/hoshikuzu/internal/services/mall/wallet"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler interface {
	QueryWallet(c *fiber.Ctx) error
	QueryWalletPage(c *fiber.Ctx) error
	UpdateWallet(c *fiber.Ctx) error
}

type WalletHandlerImpl struct {
	walletService wallet_service.WalletService
}

func NewWalletHandlerImpl(walletService wallet_service.WalletService) *WalletHandlerImpl {
	return &WalletHandlerImpl{
		walletService: walletService,
	}
}

// @Summary 查询钱包
// @Description 查询指定用户的钱包信息
// @Tags 后台管理接口/钱包
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} ent.Wallet
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/wallet/query/{user_id} [get]
func (h *WalletHandlerImpl) QueryWallet(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid user ID format",
		))
	}

	w, err := h.walletService.QueryWallet(c, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Wallet not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", w))
}

// @Summary 查询钱包分页列表
// @Description 查询所有钱包的分页列表
// @Tags 后台管理接口/钱包
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Wallet]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/wallet/page [get]
func (h *WalletHandlerImpl) QueryWalletPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, wallets, err := h.walletService.QueryWalletPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	pageResult := model.PageResult[*ent.Wallet]{
		Total:   int64(count),
		Records: wallets,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 更新钱包
// @Description 更新指定钱包的信息
// @Tags 后台管理接口/钱包
// @Accept json
// @Produce json
// @Param id path string true "钱包ID"
// @Param wallet body model.WalletUpdateReq true "钱包更新请求"
// @Success 200 {object} ent.Wallet
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/wallet/update/{id} [put]
func (h *WalletHandlerImpl) UpdateWallet(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var walletData model.WalletUpdateReq
	if err = c.BodyParser(&walletData); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedWallet, err := h.walletService.UpdateWallet(c.Context(), id, walletData)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedWallet))
}
