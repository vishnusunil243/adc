package product

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type ProductRepo struct {
	Db *gorm.DB
}

var productRepo *ProductRepo

// CreateProduct inserts a new product into the database
func (p *ProductRepo) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, *repository.RepoErr) {
	if err := p.Db.Create(&product).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return product, nil
}

// Save updates an existing product
func (p *ProductRepo) Save(ctx context.Context, product *models.Product) (*models.Product, *repository.RepoErr) {
	if err := p.Db.Save(&product).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return product, nil
}

// Delete marks products as deleted using soft delete logic
func (p *ProductRepo) Delete(ctx context.Context, Ids []string) *repository.RepoErr {
	if err := p.Db.Model(&models.Product{}).Where("id IN ?", Ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

func (p *ProductRepo) List(ctx context.Context, req *ListProductRequest) ([]*models.Product, *repository.RepoErr) {
	var products []*models.Product
	qry := p.Db.Model(&models.Product{}).Where("is_deleted = 0")

	// Filtering by IDs
	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}

	// Execute query
	if err := qry.Find(&products).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return products, nil
}

// Get fetches a product based on a provided ID or Name
func (p *ProductRepo) Get(ctx context.Context, req *GetProductRequest) (*models.Product, *repository.RepoErr) {
	var product models.Product
	qry := p.Db.Model(&models.Product{}).Where("is_deleted=0")

	// Filtering by Name or ID
	if req.Name != "" {
		qry = qry.Where("name = ?", req.Name)
	} else if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request: id or name is required", common.ErrCodeInvalidRequest)
	}

	// Execute query
	if err := qry.First(&product).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &product, nil
}

// NewProductRepo initializes a new ProductRepo instance
func NewProductRepo() *ProductRepo {
	if productRepo == nil {
		productRepo = &ProductRepo{Db: database.GetDb()}
	}
	return productRepo
}
