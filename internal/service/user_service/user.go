package user_service

import (
	"context"
	"net/http"
	"time"

	"main.go/common"
	"main.go/common/utils"
	"main.go/internal/models"
	"main.go/internal/repository/user"
	"main.go/internal/service"
)

type UserServiceApi interface {
	Signup(ctx context.Context, req *SignupReq) (*UserRes, *service.ServiceError)
	Login(ctx context.Context, req *LoginReq) (*LoginRes, *service.ServiceError)
	GetUser(ctx context.Context, req *GetUserReq) (*UserRes, *service.ServiceError)
	ListUsers(ctx context.Context, req *ListUserReq) (ListUsersRes, *service.ServiceError)
}

type UserService struct {
	userRepo *user.UserRepo
}

// Login implements UserServiceApi.
func (u *UserService) Login(ctx context.Context, req *LoginReq) (*LoginRes, *service.ServiceError) {
	user, repoErr := u.userRepo.Get(ctx, &user.GetUserRequest{
		Email: req.Email,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Invalid credentials")
	}
	valid := utils.ComparePassword(user.Password, req.Password)
	if !valid {
		return nil, service.NewServiceError("Invalid credentials", common.ErrCodeValidationFailed, http.StatusBadRequest)
	}
	token, err := utils.GenerateJWTToken(user.Id, 7*24*time.Hour)
	if err != nil {
		return nil, service.NewServiceError(err.Error(), common.ErrCodeInternal, http.StatusInternalServerError)
	}
	return NewLoginRes(user, token), nil
}

// Signup implements UserServiceApi.
func (u *UserService) Signup(ctx context.Context, req *SignupReq) (*UserRes, *service.ServiceError) {
	existingUser, repoErr := u.userRepo.Get(ctx, &user.GetUserRequest{
		Email: req.Email,
	})
	if repoErr != nil && repoErr.ErrorCode != common.ErrCodeNotFound {
		return nil, service.HandleRepoErr(repoErr, "Failed to check existing users")
	}
	if existingUser != nil {
		return nil, service.NewServiceError("User with the same email already exists", common.ErrCodeValidationFailed, http.StatusBadRequest)
	}
	encryptedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, service.NewServiceError(err.Error(), common.ErrCodeInvalidRequest, http.StatusBadRequest)
	}
	user, repoErr := u.userRepo.CreateUsers(ctx, models.NewUser(req.Name, req.Email, encryptedPassword))
	if err != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to signup")
	}
	return NewUserRes(user), nil
}

func (u *UserService) ListUsers(ctx context.Context, req *ListUserReq) (ListUsersRes, *service.ServiceError) {
	users, err := u.userRepo.List(ctx, &user.ListUserRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Ids:    req.Ids,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	return NewListUsersRes(users), nil
}

func (u *UserService) GetUser(ctx context.Context, req *GetUserReq) (*UserRes, *service.ServiceError) {
	user, err := u.userRepo.Get(ctx, &user.GetUserRequest{
		Email: req.Email,
		Id:    req.Id,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to get user")
	}
	return NewUserRes(user), nil
}

func NewUserService() UserServiceApi {
	return &UserService{
		userRepo: user.NewUserRepo(),
	}
}
