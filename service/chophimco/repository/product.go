package repository

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IProductRepo interface {
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id int) (*entity.Product, error)
	GetProductsByCategory(categoryID int) ([]entity.Product, error)
	GetProductsByBrand(brandID int) ([]entity.Product, error)
	CreateProduct(product *entity.Product) error
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id int) error

	// Variant operations
	GetVariantByID(id int) (*entity.ProductVariant, error)
	GetVariantsBySKU(sku string) (*entity.ProductVariant, error)
	CreateVariant(variant *entity.ProductVariant) error
	UpdateVariant(variant *entity.ProductVariant) error
	UpdateVariantStock(variantID int, quantity int) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Preload("Category").Preload("Brand").Preload("Variants.Switch").
		Where("is_active = ?", true).Find(&products).Error
	return products, err
}

func (r *productRepo) GetProductByID(id int) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Preload("Category").Preload("Brand").Preload("Variants.Switch").
		Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productRepo) GetProductsByCategory(categoryID int) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Preload("Category").Preload("Brand").Preload("Variants.Switch").
		Where("category_id = ? AND is_active = ?", categoryID, true).Find(&products).Error
	return products, err
}

func (r *productRepo) GetProductsByBrand(brandID int) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Preload("Category").Preload("Brand").Preload("Variants.Switch").
		Where("brand_id = ? AND is_active = ?", brandID, true).Find(&products).Error
	return products, err
}

func (r *productRepo) CreateProduct(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepo) UpdateProduct(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepo) DeleteProduct(id int) error {
	return r.db.Model(&entity.Product{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *productRepo) GetVariantByID(id int) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	err := r.db.Preload("Product").Preload("Switch").Where("id = ?", id).First(&variant).Error
	return &variant, err
}

func (r *productRepo) GetVariantsBySKU(sku string) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	err := r.db.Preload("Product").Preload("Switch").Where("sku = ?", sku).First(&variant).Error
	return &variant, err
}

func (r *productRepo) CreateVariant(variant *entity.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *productRepo) UpdateVariant(variant *entity.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *productRepo) UpdateVariantStock(variantID int, quantity int) error {
	return r.db.Model(&entity.ProductVariant{}).
		Where("id = ?", variantID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
