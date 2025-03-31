package payment

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type PaymentRepo struct {
	Db *gorm.DB
}

var paymentRepo *PaymentRepo

// CreatePayment inserts a new payment into the database
func (p *PaymentRepo) CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, *repository.RepoErr) {
	if err := p.Db.Create(&payment).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return payment, nil
}

// Save updates an existing payment
func (p *PaymentRepo) Save(ctx context.Context, payment *models.Payment) (*models.Payment, *repository.RepoErr) {
	if err := p.Db.Save(&payment).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return payment, nil
}

// Delete marks payments as deleted using soft delete logic
func (p *PaymentRepo) Delete(ctx context.Context, ids []string) *repository.RepoErr {
	if err := p.Db.Model(&models.Payment{}).Where("id IN ?", ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// List fetches a list of payments based on the request parameters
func (p *PaymentRepo) List(ctx context.Context, req *ListPaymentRequest) ([]*models.Payment, *repository.RepoErr) {
	var payments []*models.Payment
	qry := p.Db.Model(&models.Payment{}).Where("is_deleted = 0")

	// Filtering by IDs
	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}

	// Execute query
	if err := qry.Find(&payments).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return payments, nil
}

// Get fetches a payment based on a provided ID or OrderId
func (p *PaymentRepo) Get(ctx context.Context, req *GetPaymentRequest) (*models.Payment, *repository.RepoErr) {
	var payment models.Payment
	qry := p.Db.Model(&models.Payment{}).Where("is_deleted = 0")

	// Filtering by OrderId or ID
	if req.OrderId != "" {
		qry = qry.Where("order_id = ?", req.OrderId)
	} else if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request: id or order_id is required", common.ErrCodeInvalidRequest)
	}

	// Execute query
	if err := qry.First(&payment).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &payment, nil
}

// NewPaymentRepo initializes a new PaymentRepo instance
func NewPaymentRepo() *PaymentRepo {
	if paymentRepo == nil {
		paymentRepo = &PaymentRepo{Db: database.GetDb()}
	}
	return paymentRepo
}
