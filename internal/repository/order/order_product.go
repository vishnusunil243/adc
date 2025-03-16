package order

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

// OrderProductRepo handles OrderProduct CRUD operations
type OrderProductRepo struct {
	Db *gorm.DB
}

var orderProductRepo *OrderProductRepo

// CreateOrderProduct inserts a new order product into the database
func (op *OrderProductRepo) CreateOrderProduct(ctx context.Context, orderProduct *models.OrderProduct) (*models.OrderProduct, *repository.RepoErr) {
	if err := op.Db.Create(&orderProduct).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return orderProduct, nil
}

// Save updates an existing order product
func (op *OrderProductRepo) Save(ctx context.Context, orderProduct *models.OrderProduct) (*models.OrderProduct, *repository.RepoErr) {
	if err := op.Db.Save(&orderProduct).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return orderProduct, nil
}

// Delete marks order products as deleted
func (op *OrderProductRepo) Delete(ctx context.Context, Ids []string) *repository.RepoErr {
	if err := op.Db.Model(&models.OrderProduct{}).Where("id IN ?", Ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// List retrieves a list of order products
func (op *OrderProductRepo) List(ctx context.Context, req *ListOrderProductRequest) ([]*models.OrderProduct, *repository.RepoErr) {
	var orderProducts []*models.OrderProduct
	qry := op.Db.Model(&models.OrderProduct{}).Where("is_deleted = 0")

	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else if req.OrderId != "" {
		qry = qry.Where("order_id=?", req.OrderId)
	} else if len(req.OrderIds) > 0 {
		qry = qry.Where("order_id IN ?", req.OrderIds)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}

	if err := qry.Find(&orderProducts).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return orderProducts, nil
}

// Get fetches an order product based on ID
func (op *OrderProductRepo) Get(ctx context.Context, req *GetOrderProductRequest) (*models.OrderProduct, *repository.RepoErr) {
	var orderProduct models.OrderProduct
	qry := op.Db.Model(&models.OrderProduct{}).Where("is_deleted=0")

	if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request: id is required", common.ErrCodeInvalidRequest)
	}

	if err := qry.First(&orderProduct).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &orderProduct, nil
}

func (op *OrderProductRepo) BulkCreate(ctx context.Context, orderProducts []*models.OrderProduct) *repository.RepoErr {
	if err := op.Db.WithContext(ctx).Create(&orderProducts).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// NewOrderProductRepo initializes a new OrderProductRepo instance
func NewOrderProductRepo() *OrderProductRepo {
	if orderProductRepo == nil {
		orderProductRepo = &OrderProductRepo{Db: database.GetDb()}
	}
	return orderProductRepo
}
