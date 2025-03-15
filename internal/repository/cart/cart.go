package cart

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type CartRepo struct {
	Db *gorm.DB
}

var cartRepo *CartRepo

// AddToCart inserts a new cart entry into the database
func (c *CartRepo) AddToCart(ctx context.Context, cart *models.Cart) (*models.Cart, *repository.RepoErr) {
	if err := c.Db.Create(&cart).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return cart, nil
}

// UpdateCart updates an existing cart entry
func (c *CartRepo) UpdateCart(ctx context.Context, id string, cart *models.Cart) (*models.Cart, *repository.RepoErr) {
	if err := c.Db.Where("id=?", id).Save(&cart).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return cart, nil
}

// DeleteCart marks cart items as deleted using soft delete logic
func (c *CartRepo) DeleteCart(ctx context.Context, Ids []string) *repository.RepoErr {
	if err := c.Db.Model(&models.Cart{}).Where("id IN ?", Ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// ListCart retrieves all cart items
func (c *CartRepo) ListCart(ctx context.Context, req *ListCartRequest) ([]*models.Cart, *repository.RepoErr) {
	var carts []*models.Cart
	qry := c.Db.Model(&models.Cart{}).Where("is_deleted = 0")

	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}

	if req.UserId != "" {
		qry = qry.Where("user_id=?", req.UserId)
	}

	if err := qry.Find(&carts).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return carts, nil
}

// GetCart retrieves a cart entry by ID
func (c *CartRepo) GetCart(ctx context.Context, req *GetCartRequest) (*models.Cart, *repository.RepoErr) {
	var cart models.Cart
	qry := c.Db.Model(&models.Cart{}).Where("is_deleted=0")

	if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else if req.UserId != "" && req.ProductId != "" {
		qry = qry.Where("user_id=? AND product_id=?", req.UserId, req.ProductId)
	} else {
		return nil, repository.NewRepoErr("invalid request: id is required", common.ErrCodeInvalidRequest)
	}

	if err := qry.First(&cart).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &cart, nil
}

// NewCartRepo initializes a new CartRepo instance
func NewCartRepo() *CartRepo {
	if cartRepo == nil {
		cartRepo = &CartRepo{Db: database.GetDb()}
	}
	return cartRepo
}
