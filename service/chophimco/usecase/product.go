package usecase

import (
	"context"
	"errors"

	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/response"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
	"gorm.io/gorm"
)

type IProductUsecase interface {
	GetAllProducts(ctx context.Context) ([]response.ProductResponse, error)
	GetProductByID(ctx context.Context, id int) (*response.ProductResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID int) ([]response.ProductResponse, error)
	GetProductsByBrand(ctx context.Context, brandID int) ([]response.ProductResponse, error)
	CreateProduct(ctx context.Context, req request.CreateProduct) error
	UpdateProduct(ctx context.Context, req request.UpdateProduct) error
	DeleteProduct(ctx context.Context, id int) error
	CreateProductVariant(ctx context.Context, req request.CreateProductVariant) error
	UpdateProductVariant(ctx context.Context, req request.UpdateProductVariant) error
}

type productUsecase struct {
	repo repository.IProductRepo
}

func NewProductUsecase(repo repository.IProductRepo) IProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]response.ProductResponse, error) {
	products, err := u.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return u.mapProductsToResponse(products), nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id int) (*response.ProductResponse, error) {
	product, err := u.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	resp := u.mapProductToResponse(product)
	return &resp, nil
}

func (u *productUsecase) GetProductsByCategory(ctx context.Context, categoryID int) ([]response.ProductResponse, error) {
	products, err := u.repo.GetProductsByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return u.mapProductsToResponse(products), nil
}

func (u *productUsecase) GetProductsByBrand(ctx context.Context, brandID int) ([]response.ProductResponse, error) {
	products, err := u.repo.GetProductsByBrand(brandID)
	if err != nil {
		return nil, err
	}
	return u.mapProductsToResponse(products), nil
}

func (u *productUsecase) CreateProduct(ctx context.Context, req request.CreateProduct) error {
	product := &entity.Product{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		Description: req.Description,
		BasePrice:   req.BasePrice,
		IsActive:    true,
	}
	return u.repo.CreateProduct(product)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, req request.UpdateProduct) error {
	product, err := u.repo.GetProductByID(req.ID)
	if err != nil {
		return err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.BrandID != nil {
		product.BrandID = req.BrandID
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.BasePrice > 0 {
		product.BasePrice = req.BasePrice
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	return u.repo.UpdateProduct(product)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	return u.repo.DeleteProduct(id)
}

func (u *productUsecase) CreateProductVariant(ctx context.Context, req request.CreateProductVariant) error {
	// Check if SKU exists
	_, err := u.repo.GetVariantsBySKU(req.SKU)
	if err == nil {
		return errors.New("SKU already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	variant := &entity.ProductVariant{
		ProductID:      req.ProductID,
		SwitchID:       req.SwitchID,
		Layout:         req.Layout,
		ConnectionType: req.ConnectionType,
		Hotswap:        req.Hotswap,
		LedType:        req.LedType,
		Price:          req.Price,
		Stock:          req.Stock,
		SKU:            req.SKU,
	}
	return u.repo.CreateVariant(variant)
}

func (u *productUsecase) UpdateProductVariant(ctx context.Context, req request.UpdateProductVariant) error {
	variant, err := u.repo.GetVariantByID(req.ID)
	if err != nil {
		return err
	}

	if req.SwitchID != nil {
		variant.SwitchID = req.SwitchID
	}
	if req.Layout != "" {
		variant.Layout = req.Layout
	}
	if req.ConnectionType != "" {
		variant.ConnectionType = req.ConnectionType
	}
	if req.Hotswap != nil {
		variant.Hotswap = *req.Hotswap
	}
	if req.LedType != "" {
		variant.LedType = req.LedType
	}
	if req.Price > 0 {
		variant.Price = req.Price
	}
	if req.Stock != nil {
		variant.Stock = *req.Stock
	}

	return u.repo.UpdateVariant(variant)
}

func (u *productUsecase) mapProductsToResponse(products []entity.Product) []response.ProductResponse {
	result := make([]response.ProductResponse, 0, len(products))
	for _, p := range products {
		result = append(result, u.mapProductToResponse(&p))
	}
	return result
}

func (u *productUsecase) mapProductToResponse(product *entity.Product) response.ProductResponse {
	resp := response.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		BasePrice:   product.BasePrice,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
	}

	if product.Category != nil {
		categoryName := product.Category.Name
		resp.Category = &categoryName
	}

	if product.Brand != nil {
		brandName := product.Brand.Name
		resp.Brand = &brandName
	}

	if len(product.Variants) > 0 {
		variants := make([]response.ProductVariantResponse, 0, len(product.Variants))
		for _, v := range product.Variants {
			variantResp := response.ProductVariantResponse{
				ID:             v.ID,
				ProductID:      v.ProductID,
				Layout:         v.Layout,
				ConnectionType: v.ConnectionType,
				Hotswap:        v.Hotswap,
				LedType:        v.LedType,
				Price:          v.Price,
				Stock:          v.Stock,
				SKU:            v.SKU,
			}
			if v.Switch != nil {
				switchName := v.Switch.Name
				variantResp.Switch = &switchName
			}
			variants = append(variants, variantResp)
		}
		resp.Variants = variants
	}

	return resp
}
