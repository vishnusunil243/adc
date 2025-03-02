package common

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NewBaseResponse(success bool, data interface{}) *BaseResponse {
	return &BaseResponse{
		Success: success,
		Data:    data,
	}
}

func NewResponse(ctx echo.Context, success bool, data interface{}, statusCode int) error {
	return ctx.JSON(statusCode, NewBaseResponse(success, data))
}
