package user_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/user_service"
)

type UserHandler struct {
	userService user_service.UserServiceApi
}

func (u *UserHandler) Signup(c echo.Context) error {
	var req user_service.SignupReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	userRes, err := u.userService.Signup(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, userRes, http.StatusCreated)
}

func (u *UserHandler) Login(c echo.Context) error {
	var req user_service.LoginReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	userRes, err := u.userService.Login(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, userRes, http.StatusCreated)
}

func (u *UserHandler) List(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	// Convert them to integers (default to 10 and 0 if not provided)
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)
	req := user_service.ListUserReq{
		Limit:  limitInt,
		Offset: offsetInt,
	}
	listRes, err := u.userService.ListUsers(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, listRes, http.StatusOK)
}

func (u *UserHandler) Get(c echo.Context) error {
	id := c.Param("id")
	userRes, err := u.userService.GetUser(c.Request().Context(), &user_service.GetUserReq{
		Id: id,
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, userRes, http.StatusOK)
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: user_service.NewUserService(),
	}
}
