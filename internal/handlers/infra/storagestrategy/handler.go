package storagestrategy

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/internal/infra/storage"
	storagestrategy_service "github.com/shuTwT/hoshikuzu/internal/services/infra/storagestrategy"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type StorageStrategyHandler struct {
	storageStrategyService storagestrategy_service.StorageStrategyService
}

func NewStorageStrategyHandler(storageStrategyService storagestrategy_service.StorageStrategyService) *StorageStrategyHandler {
	return &StorageStrategyHandler{
		storageStrategyService: storageStrategyService,
	}
}

// @Summary 获取存储策略分页
// @Description 获取所有存储策略的分页
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "策略名称"
// @Param type query string false "存储类型"
// @Param master query bool false "是否默认"
// @Success 200 {object} model.HttpSuccess{data=[]ent.StorageStrategy}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/page [get]
func (h *StorageStrategyHandler) ListStorageStrategyPage(c *fiber.Ctx) error {
	var pageReq model.StorageStrategyPageReq
	if err := c.QueryParser(&pageReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, strategies, err := h.storageStrategyService.ListStorageStrategyPageWithQuery(c.Context(), pageReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resps := make([]model.StorageStrategyResp, 0, len(strategies))
	for _, strategy := range strategies {
		resps = append(resps, model.StorageStrategyResp{
			ID:        strategy.ID,
			Name:      strategy.Name,
			Type:      string(strategy.Type),
			Master:    strategy.Master,
			NodeID:    strategy.NodeID,
			Endpoint:  strategy.Endpoint,
			Region:    strategy.Region,
			Bucket:    strategy.Bucket,
			BasePath:  strategy.BasePath,
			Domain:    strategy.Domain,
			AccessKey: strategy.AccessKey,
			SecretKey: strategy.SecretKey,
			CreatedAt: model.LocalTime(strategy.CreatedAt),
		})
	}

	pageResult := model.PageResult[model.StorageStrategyResp]{
		Total:   int64(count),
		Records: resps,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 获取所有存储策略列表
// @Description 获取所有存储策略的列表，包括默认策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.StorageStrategyListResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/list [get]
func (h *StorageStrategyHandler) ListStorageStrategy(c *fiber.Ctx) error {
	strategies, err := h.storageStrategyService.ListStorageStrategy(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	var strategyList []model.StorageStrategyListResp
	for _, strategy := range strategies {
		strategyList = append(strategyList, model.StorageStrategyListResp{
			ID:     strategy.ID,
			Name:   strategy.Name,
			Type:   string(strategy.Type),
			Master: strategy.Master,
		})
	}
	return c.JSON(model.NewSuccess("success", strategyList))
}

// @Summary 创建存储策略
// @Description 创建一个新的存储策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param strategy body model.StorageStrategyCreateReq true "Storage Strategy Create Request"
// @Success 200 {object} model.HttpSuccess{data=ent.StorageStrategy}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/create [post]
func (h *StorageStrategyHandler) CreateStorageStrategy(c *fiber.Ctx) error {
	var strategy *model.StorageStrategyCreateReq
	if err := c.BodyParser(&strategy); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newStrategy, err := h.storageStrategyService.CreateStorageStrategy(c.Context(), strategy.Name, strategy.Type, strategy.NodeID, strategy.Endpoint, strategy.Region, strategy.Bucket, strategy.AccessKey, strategy.SecretKey, strategy.BasePath, strategy.Domain, strategy.Master)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newStrategy))
}

// @Summary 更新存储策略
// @Description 更新指定ID的存储策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param id path int true "Storage Strategy ID"
// @Param strategy body model.StorageStrategyUpdateReq true "Storage Strategy Update Request"
// @Success 200 {object} model.HttpSuccess{data=ent.StorageStrategy}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/update/{id} [put]
func (h *StorageStrategyHandler) UpdateStorageStrategy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var strategy *model.StorageStrategyUpdateReq
	if err = c.BodyParser(&strategy); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newStrategy, err := h.storageStrategyService.UpdateStorageStrategy(c.Context(), id, strategy.Name, strategy.Type, strategy.NodeID, strategy.Endpoint, strategy.Region, strategy.Bucket, strategy.AccessKey, strategy.SecretKey, strategy.BasePath, strategy.Domain, strategy.Master)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	storage.ClearCache()
	return c.JSON(model.NewSuccess("success", newStrategy))
}

// @Summary 查询存储策略
// @Description 查询指定ID的存储策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param id path int true "Storage Strategy ID"
// @Success 200 {object} model.HttpSuccess{data=ent.StorageStrategy}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/query/{id} [get]
func (h *StorageStrategyHandler) QueryStorageStrategy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	strategy, err := h.storageStrategyService.QueryStorageStrategy(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", strategy))
}

// @Summary 删除存储策略
// @Description 删除指定ID的存储策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param id path int true "Storage Strategy ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/delete/{id} [delete]
func (h *StorageStrategyHandler) DeleteStorageStrategy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.storageStrategyService.DeleteStorageStrategy(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	storage.ClearCache()
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 设置默认存储策略
// @Description 设置指定ID的存储策略为默认策略
// @Tags 后台管理接口/存储策略
// @Accept json
// @Produce json
// @Param id path int true "Storage Strategy ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/storage-strategy/default/{id} [put]
func (h *StorageStrategyHandler) SetDefaultStorageStrategy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.storageStrategyService.SetDefaultStorageStrategy(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	storage.ClearCache()
	return c.JSON(model.NewSuccess("success", nil))
}
