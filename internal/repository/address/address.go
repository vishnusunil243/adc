package address

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/database"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type AddressRepo struct {
	Db *gorm.DB
}

var addressRepo *AddressRepo

// CreateAddress inserts a new address into the database
func (a *AddressRepo) CreateAddress(ctx context.Context, address *models.Address) (*models.Address, *repository.RepoErr) {
	if err := a.Db.Create(&address).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return address, nil
}

// Save updates an existing address
func (a *AddressRepo) Save(ctx context.Context, address *models.Address) (*models.Address, *repository.RepoErr) {
	if err := a.Db.Where("id=?", address.Id).Save(&address).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return address, nil
}

// Delete marks addresses as deleted using soft delete logic
func (a *AddressRepo) Delete(ctx context.Context, ids []string) *repository.RepoErr {
	if err := a.Db.Model(&models.Address{}).Where("id IN ?", ids).Updates(common.GetFieldsForDelete()).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

// List fetches a list of addresses based on the request filters
func (a *AddressRepo) List(ctx context.Context, req *ListAddressRequest) ([]*models.Address, *repository.RepoErr) {
	var addresses []*models.Address
	qry := a.Db.Model(&models.Address{}).Where("is_deleted = 0")

	// Filtering by IDs
	if req.Ids != nil {
		qry = qry.Where("id IN ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}
	if req.UserId != "" {
		qry = qry.Where("user_id = ?", req.UserId)
	}

	// Execute query
	if err := qry.Find(&addresses).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return addresses, nil
}

// Get fetches an address based on a provided ID or UserId
func (a *AddressRepo) Get(ctx context.Context, req *GetAddressRequest) (*models.Address, *repository.RepoErr) {
	var address models.Address
	qry := a.Db.Model(&models.Address{}).Where("is_deleted=0")

	// Filtering by UserId or Id
	if req.Id != "" {
		qry = qry.Where("id = ?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request: id or user_id is required", common.ErrCodeInvalidRequest)
	}

	// Execute query
	if err := qry.First(&address).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return &address, nil
}

// NewAddressRepo initializes a new AddressRepo instance
func NewAddressRepo() *AddressRepo {
	if addressRepo == nil {
		addressRepo = &AddressRepo{Db: database.GetDb()}
	}
	return addressRepo
}
