package order

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type OrderRepo struct {
	Db *gorm.DB
}

var orderRepo *OrderRepo

// CreateOrder inserts a new order into the database
func (o *OrderRepo) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, *repository.RepoErr) {
	if err := o.Db.Create(&order).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return order, nil
}

// Save updates an existing order
func (o *OrderRepo) Save(ctx context.Context, order *models.Order) (*models.Order, *repository.RepoErr) {
	if err := o.Db.Where("id=?", order.Id).Save(&order).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return order, nil
}

// Delete marks orders as deleted using soft delete logic
func (o *OrderRepo) Delete(ctx context.Context, Ids []string) *repository.RepoErr {
	if err := o.Db.Model(&models.Order{}).Where("id IN ?", Ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// List retrieves a list of orders
func (o *OrderRepo) List(ctx context.Context, req *ListOrderRequest) ([]*models.Order, *repository.RepoErr) {
	var orders []*models.Order
	qry := o.Db.Model(&models.Order{}).Where("is_deleted = 0")

	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}
	if req.UserId != "" {
		qry = qry.Where("user_id=?", req.UserId)
	}

	if err := qry.Find(&orders).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return orders, nil
}

// Get fetches an order based on provided ID
func (o *OrderRepo) Get(ctx context.Context, req *GetOrderRequest) (*models.Order, *repository.RepoErr) {
	var order models.Order
	qry := o.Db.Model(&models.Order{}).Where("is_deleted=0")

	if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request: id is required", common.ErrCodeInvalidRequest)
	}

	if err := qry.First(&order).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &order, nil
}

// NewOrderRepo initializes a new OrderRepo instance
func NewOrderRepo() *OrderRepo {
	if orderRepo == nil {
		orderRepo = &OrderRepo{Db: database.GetDb()}
	}
	return orderRepo
}
