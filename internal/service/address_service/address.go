package address_service

import (
	"context"

	"main.go/common/utils"
	"main.go/internal/models"
	"main.go/internal/repository/address"
	"main.go/internal/repository/user"
	"main.go/internal/service"
)

type AddressServiceApi interface {
	CreateAddress(ctx context.Context, req *CreateAddressReq) (*AddressRes, *service.ServiceError)
	UpdateAddress(ctx context.Context, req *UpdateAddressReq) (*AddressRes, *service.ServiceError)
	DeleteAddress(ctx context.Context, req *DeleteAddressReq) *service.ServiceError
	GetAddress(ctx context.Context, req *GetAddressReq) (*AddressRes, *service.ServiceError)
	ListAddresses(ctx context.Context, req *ListAddressReq) (ListAddressRes, *service.ServiceError)
}

type AddressService struct {
	addressRepo *address.AddressRepo
	userRepo    *user.UserRepo
}

// CreateAddress implements AddressServiceApi.
func (a *AddressService) CreateAddress(ctx context.Context, req *CreateAddressReq) (*AddressRes, *service.ServiceError) {
	address, repoErr := a.addressRepo.CreateAddress(ctx, models.NewAddress(req.UserId, req.Country, req.State, req.City, req.Pincode, req.Street, req.Area, req.CreatedBy))
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to create address")
	}
	return a.getAddressResponse(ctx, address)
}

// UpdateAddress implements AddressServiceApi.
func (a *AddressService) UpdateAddress(ctx context.Context, req *UpdateAddressReq) (*AddressRes, *service.ServiceError) {
	existingAddress, repoErr := a.addressRepo.Get(ctx, &address.GetAddressRequest{Id: req.Id})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Address not found")
	}

	// Update Address fields as necessary
	existingAddress.UpdateDetails(req.Country, req.State, req.City, req.Pincode, req.Street, req.Area)

	updatedAddress, repoErr := a.addressRepo.Save(ctx, existingAddress)
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to update address")
	}
	return a.getAddressResponse(ctx, updatedAddress)
}

// DeleteAddress implements AddressServiceApi.
func (a *AddressService) DeleteAddress(ctx context.Context, req *DeleteAddressReq) *service.ServiceError {
	repoErr := a.addressRepo.Delete(ctx, req.Ids)
	if repoErr != nil {
		return service.HandleRepoErr(repoErr, "Failed to delete address")
	}
	return nil
}

// GetAddress implements AddressServiceApi.
func (a *AddressService) GetAddress(ctx context.Context, req *GetAddressReq) (*AddressRes, *service.ServiceError) {
	address, repoErr := a.addressRepo.Get(ctx, &address.GetAddressRequest{Id: req.Id})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to get address")
	}
	return a.getAddressResponse(ctx, address)
}

// ListAddresses implements AddressServiceApi.
func (a *AddressService) ListAddresses(ctx context.Context, req *ListAddressReq) (ListAddressRes, *service.ServiceError) {
	addresses, repoErr := a.addressRepo.List(ctx, &address.ListAddressRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		UserId: utils.GetCurrentUser(ctx),
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to list addresses")
	}
	return a.getListAddressRes(ctx, addresses)
}

func (a *AddressService) getAddressResponse(ctx context.Context, address *models.Address) (*AddressRes, *service.ServiceError) {
	users, err := a.userRepo.List(ctx, &user.ListUserRequest{
		Ids: address.GetAuditFieldsUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	return NewAddressRes(address, users), nil
}

func (a *AddressService) getListAddressRes(ctx context.Context, addresses models.ListAddress) (ListAddressRes, *service.ServiceError) {
	users, err := a.userRepo.List(ctx, &user.ListUserRequest{
		Ids: addresses.GetUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	return NewListAddressRes(addresses, users), nil
}

// NewAddressService initializes a new AddressService instance
func NewAddressService() AddressServiceApi {
	return &AddressService{
		addressRepo: address.NewAddressRepo(),
		userRepo:    user.NewUserRepo(),
	}
}
