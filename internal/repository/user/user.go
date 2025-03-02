package user

import (
	"context"

	"gorm.io/gorm"
	"main.go/common"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type UserRepo struct {
	Db *gorm.DB
}

var userRepo *UserRepo

func (u *UserRepo) CreateUsers(ctx context.Context, user *models.User) (*models.User, *repository.RepoErr) {
	if err := u.Db.Create(&user).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return user, nil
}

func (u *UserRepo) Save(ctx context.Context, user *models.User) (*models.User, *repository.RepoErr) {
	if err := u.Db.Save(&user).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return user, nil
}

func (u *UserRepo) Delete(ctx context.Context, Ids []string) *repository.RepoErr {
	if err := u.Db.Updates(common.GetFieldsForDelete()).Where("id in ?", Ids).Error; err != nil {
		return repository.HandleDBError(err)
	}
	return nil
}

func (u *UserRepo) List(ctx context.Context, req *ListUserRequest) ([]*models.User, *repository.RepoErr) {
	var users []*models.User
	qry := u.Db.Model(&models.User{}).Where("is_deleted=0")
	if req.Ids != nil {
		qry = qry.Where("id in ?", req.Ids)
	} else {
		qry = qry.Limit(req.Limit).Offset(req.Offset)
	}
	if err := qry.Find(&users).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return users, nil
}

func (u *UserRepo) Get(ctx context.Context, req *GetUserRequest) (*models.User, *repository.RepoErr) {
	var user *models.User
	qry := u.Db.Model(&models.User{}).Where("is_deleted=0")
	if req.Email != "" {
		qry = qry.Where("email=?", req.Email)
	} else if req.Name != "" {
		qry = qry.Where("name=?", req.Name)
	} else if req.Id != "" {
		qry = qry.Where("id =?", req.Id)
	} else {
		return nil, repository.NewRepoErr("invalid request id or email is required", common.ErrCodeInvalidRequest)
	}
	if err := qry.First(&user).Error; err != nil {
		return nil, repository.HandleDBError(err)
	}
	return user, nil

}

func NewUserRepo() *UserRepo {
	if userRepo == nil {
		userRepo = &UserRepo{}
	}
	return userRepo
}
