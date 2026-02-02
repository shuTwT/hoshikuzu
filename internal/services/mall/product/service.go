package product

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/product"
	"github.com/shuTwT/hoshikuzu/pkg/cache"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *model.ProductCreateReq) (*ent.Product, error)
	UpdateProduct(ctx context.Context, id int, req *model.ProductUpdateReq) (*ent.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	GetProduct(ctx context.Context, id int) (*ent.Product, error)
	ListProducts(ctx context.Context, page, pageSize int) ([]*ent.Product, int, error)
	ListAllProducts(ctx context.Context) ([]*ent.Product, error)
	BatchUpdateProducts(ctx context.Context, ids []int, req *model.ProductBatchUpdateReq) error
	BatchDeleteProducts(ctx context.Context, ids []int) error
	SearchProducts(ctx context.Context, req model.ProductSearchReq) ([]*model.ProductSearchResp, int, error)
}

type ProductServiceImpl struct {
	client *ent.Client
}

func NewProductServiceImpl(client *ent.Client) *ProductServiceImpl {
	return &ProductServiceImpl{client: client}
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *model.ProductCreateReq) (*ent.Product, error) {
	createBuilder := s.client.Product.Create().
		SetName(req.Name).
		SetSku(req.Sku).
		SetPrice(req.Price).
		SetStock(req.Stock)

	if req.Description != nil {
		createBuilder.SetDescription(*req.Description)
	}
	if req.ShortDescription != nil {
		createBuilder.SetShortDescription(*req.ShortDescription)
	}
	if req.OriginalPrice != nil {
		createBuilder.SetOriginalPrice(*req.OriginalPrice)
	}
	if req.CostPrice != nil {
		createBuilder.SetCostPrice(*req.CostPrice)
	}
	if req.MinStock != nil {
		createBuilder.SetMinStock(*req.MinStock)
	}
	if req.CategoryID != nil {
		createBuilder.SetCategoryID(*req.CategoryID)
	}
	if req.Brand != nil {
		createBuilder.SetBrand(*req.Brand)
	}
	if req.Unit != nil {
		createBuilder.SetUnit(*req.Unit)
	}
	if req.Weight != nil {
		createBuilder.SetWeight(*req.Weight)
	}
	if req.Volume != nil {
		createBuilder.SetVolume(*req.Volume)
	}
	if req.Images != nil {
		createBuilder.SetImages(req.Images)
	}
	if req.Attributes != nil {
		createBuilder.SetAttributes(req.Attributes)
	}
	if req.Tags != nil {
		createBuilder.SetTags(req.Tags)
	}
	if req.Active != nil {
		createBuilder.SetActive(*req.Active)
	}
	if req.Featured != nil {
		createBuilder.SetFeatured(*req.Featured)
	}
	if req.Digital != nil {
		createBuilder.SetDigital(*req.Digital)
	}
	if req.MetaTitle != nil {
		createBuilder.SetMetaTitle(*req.MetaTitle)
	}
	if req.MetaDescription != nil {
		createBuilder.SetMetaDescription(*req.MetaDescription)
	}
	if req.MetaKeywords != nil {
		createBuilder.SetMetaKeywords(*req.MetaKeywords)
	}
	if req.SortOrder != nil {
		createBuilder.SetSortOrder(*req.SortOrder)
	}

	return createBuilder.Save(ctx)
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, id int, req *model.ProductUpdateReq) (*ent.Product, error) {
	updateBuilder := s.client.Product.UpdateOneID(id)

	if req.Name != nil {
		updateBuilder.SetName(*req.Name)
	}
	if req.Description != nil {
		updateBuilder.SetDescription(*req.Description)
	}
	if req.ShortDescription != nil {
		updateBuilder.SetShortDescription(*req.ShortDescription)
	}
	if req.Sku != nil {
		updateBuilder.SetSku(*req.Sku)
	}
	if req.Price != nil {
		updateBuilder.SetPrice(*req.Price)
	}
	if req.OriginalPrice != nil {
		updateBuilder.SetOriginalPrice(*req.OriginalPrice)
	}
	if req.CostPrice != nil {
		updateBuilder.SetCostPrice(*req.CostPrice)
	}
	if req.Stock != nil {
		updateBuilder.SetStock(*req.Stock)
	}
	if req.MinStock != nil {
		updateBuilder.SetMinStock(*req.MinStock)
	}
	if req.Sales != nil {
		updateBuilder.SetSales(*req.Sales)
	}
	if req.CategoryID != nil {
		updateBuilder.SetCategoryID(*req.CategoryID)
	}
	if req.Brand != nil {
		updateBuilder.SetBrand(*req.Brand)
	}
	if req.Unit != nil {
		updateBuilder.SetUnit(*req.Unit)
	}
	if req.Weight != nil {
		updateBuilder.SetWeight(*req.Weight)
	}
	if req.Volume != nil {
		updateBuilder.SetVolume(*req.Volume)
	}
	if req.Images != nil {
		updateBuilder.SetImages(req.Images)
	}
	if req.Attributes != nil {
		updateBuilder.SetAttributes(req.Attributes)
	}
	if req.Tags != nil {
		updateBuilder.SetTags(req.Tags)
	}
	if req.Active != nil {
		updateBuilder.SetActive(*req.Active)
	}
	if req.Featured != nil {
		updateBuilder.SetFeatured(*req.Featured)
	}
	if req.Digital != nil {
		updateBuilder.SetDigital(*req.Digital)
	}
	if req.MetaTitle != nil {
		updateBuilder.SetMetaTitle(*req.MetaTitle)
	}
	if req.MetaDescription != nil {
		updateBuilder.SetMetaDescription(*req.MetaDescription)
	}
	if req.MetaKeywords != nil {
		updateBuilder.SetMetaKeywords(*req.MetaKeywords)
	}
	if req.SortOrder != nil {
		updateBuilder.SetSortOrder(*req.SortOrder)
	}

	return updateBuilder.Save(ctx)
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id int) error {
	return s.client.Product.DeleteOneID(id).Exec(ctx)
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, id int) (*ent.Product, error) {
	return s.client.Product.Query().
		Where(product.ID(id)).
		First(ctx)
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, page, pageSize int) ([]*ent.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total, err := s.client.Product.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	products, err := s.client.Product.Query().
		Order(ent.Desc(product.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)

	return products, total, err
}

func (s *ProductServiceImpl) ListAllProducts(ctx context.Context) ([]*ent.Product, error) {
	return s.client.Product.Query().
		Order(ent.Desc(product.FieldCreatedAt)).
		All(ctx)
}

func (s *ProductServiceImpl) BatchUpdateProducts(ctx context.Context, ids []int, req *model.ProductBatchUpdateReq) error {
	if len(ids) == 0 {
		return nil
	}

	updateBuilder := s.client.Product.Update().
		Where(product.IDIn(ids...))

	if req.Active != nil {
		updateBuilder.SetActive(*req.Active)
	}

	_, err := updateBuilder.Save(ctx)
	return err
}

func (s *ProductServiceImpl) BatchDeleteProducts(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.client.Product.Delete().
		Where(product.IDIn(ids...)).
		Exec(ctx)
	return err
}

func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req model.ProductSearchReq) ([]*model.ProductSearchResp, int, error) {
	cacheKey := fmt.Sprintf("product:search:%s:%d:%d", req.Keyword, req.Page, req.Size)

	if cached, found := cache.GetCache().Get(cacheKey); found {
		if result, ok := cached.([]*model.ProductSearchResp); ok {
			return result, len(result), nil
		}
	}

	keyword := strings.ToLower(req.Keyword)

	products, err := s.client.Product.Query().
		Where(
			product.Or(
				product.NameContains(keyword),
				product.SkuContains(keyword),
				product.DescriptionContains(keyword),
				product.ShortDescriptionContains(keyword),
				product.BrandContains(keyword),
			),
		).
		Order(ent.Desc(product.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	var results []*model.ProductSearchResp

	for _, p := range products {
		relevance := s.calculateProductRelevance(p, keyword)
		if relevance > 0 {
			results = append(results, &model.ProductSearchResp{
				ID:               p.ID,
				Name:             p.Name,
				ShortDescription: &p.ShortDescription,
				Sku:              p.Sku,
				Price:            p.Price,
				OriginalPrice:    &p.OriginalPrice,
				Stock:            p.Stock,
				Sales:            p.Sales,
				Brand:            &p.Brand,
				Images:           p.Images,
				Active:           p.Active,
				Relevance:        relevance,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Relevance == results[j].Relevance {
			return results[i].Sales > results[j].Sales
		}
		return results[i].Relevance > results[j].Relevance
	})

	total := len(results)

	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start >= total {
		return []*model.ProductSearchResp{}, total, nil
	}
	if end > total {
		end = total
	}

	pagedResults := results[start:end]

	cache.GetCache().Set(cacheKey, pagedResults, 5*time.Minute)

	return pagedResults, total, nil
}

func (s *ProductServiceImpl) calculateProductRelevance(p *ent.Product, keyword string) float64 {
	var relevance float64 = 0

	name := strings.ToLower(p.Name)
	sku := strings.ToLower(p.Sku)
	description := strings.ToLower(p.Description)
	shortDescription := strings.ToLower(p.ShortDescription)
	brand := strings.ToLower(p.Brand)

	if strings.Contains(name, keyword) {
		if name == keyword {
			relevance += 10.0
		} else if strings.HasPrefix(name, keyword) {
			relevance += 8.0
		} else {
			relevance += 5.0
		}
	}

	if strings.Contains(sku, keyword) {
		if sku == keyword {
			relevance += 8.0
		} else {
			relevance += 5.0
		}
	}

	if strings.Contains(shortDescription, keyword) {
		relevance += 3.0
	}

	if strings.Contains(description, keyword) {
		relevance += 2.0
	}

	if strings.Contains(brand, keyword) {
		relevance += 1.0
	}

	return relevance
}
