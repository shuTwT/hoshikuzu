package product

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	product_service "github.com/shuTwT/hoshikuzu/internal/services/mall/product"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler interface {
	ListProducts(c *fiber.Ctx) error
	ListProductsPage(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	QueryProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	BatchUpdateProducts(c *fiber.Ctx) error
	BatchDeleteProducts(c *fiber.Ctx) error
	SearchProducts(c *fiber.Ctx) error
}

type ProductHandlerImpl struct {
	productService product_service.ProductService
}

func NewProductHandlerImpl(productService product_service.ProductService) *ProductHandlerImpl {
	return &ProductHandlerImpl{
		productService: productService,
	}
}

// @Summary 查询所有商品
// @Description 查询所有商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Product}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/list [get]
func (h *ProductHandlerImpl) ListProducts(c *fiber.Ctx) error {
	products, err := h.productService.ListAllProducts(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", products))
}

// @Summary 分页查询商品
// @Description 分页查询商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Product]}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/page [get]
func (h *ProductHandlerImpl) ListProductsPage(c *fiber.Ctx) error {
	var pageQuery model.PageQuery

	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if pageQuery.Page < 1 {
		pageQuery.Page = 1
	}
	if pageQuery.Size < 1 {
		pageQuery.Size = 10
	}

	products, total, err := h.productService.ListProducts(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.Product]{
		Total:   int64(total),
		Records: products,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建商品
// @Description 创建一个新商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param product body model.ProductCreateReq true "商品创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Product}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/create [post]
func (h *ProductHandlerImpl) CreateProduct(c *fiber.Ctx) error {
	var req *model.ProductCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	product, err := h.productService.CreateProduct(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", product))
}

// @Summary 更新商品
// @Description 更新指定商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param product body model.ProductUpdateReq true "商品更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Product}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/update/{id} [put]
func (h *ProductHandlerImpl) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req *model.ProductUpdateReq
	if err = c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	product, err := h.productService.UpdateProduct(c.Context(), id, req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", product))
}

// @Summary 查询商品
// @Description 查询指定商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Product}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/query/{id} [get]
func (h *ProductHandlerImpl) QueryProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	product, err := h.productService.GetProduct(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", product))
}

// @Summary 删除商品
// @Description 删除指定商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/delete/{id} [delete]
func (h *ProductHandlerImpl) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	if err := h.productService.DeleteProduct(c.Context(), id); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量更新商品
// @Description 批量更新指定商品的信息
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param products body model.ProductBatchUpdateReq true "商品批量更新请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/batch-update [put]
func (h *ProductHandlerImpl) BatchUpdateProducts(c *fiber.Ctx) error {
	var req *model.ProductBatchUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.productService.BatchUpdateProducts(c.Context(), req.IDs, req); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量删除商品
// @Description 批量删除指定商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param products body model.ProductBatchDeleteReq true "商品批量删除请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/batch-delete [delete]
func (h *ProductHandlerImpl) BatchDeleteProducts(c *fiber.Ctx) error {
	var req *model.ProductBatchDeleteReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.productService.BatchDeleteProducts(c.Context(), req.IDs); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 搜索商品
// @Description 根据搜索条件查询商品
// @Tags 后台管理接口/商品
// @Accept json
// @Produce json
// @Param name query string false "商品名称"
// @Param category_id query int false "商品分类ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.ProductSearchResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/product/search [get]
func (h *ProductHandlerImpl) SearchProducts(c *fiber.Ctx) error {
	var req model.ProductSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	results, total, err := h.productService.SearchProducts(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*model.ProductSearchResp]{
		Total:   int64(total),
		Records: results,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}
