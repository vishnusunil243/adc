package order_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/order_service"
)

type OrderHandler struct {
	orderService order_service.OrderServiceApi
}

func (oh *OrderHandler) AddOrder(c echo.Context) error {
	var req order_service.AddOrderRequest
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	if err := oh.orderService.AddOrder(c.Request().Context(), &req); err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, "Order created successfully", http.StatusCreated)
}

func (oh *OrderHandler) ListOrders(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	req := order_service.ListOrderRequest{
		Limit:  limitInt,
		Offset: offsetInt,
	}
	listRes, err := oh.orderService.ListOrders(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, listRes, http.StatusOK)
}

func (oh *OrderHandler) GetOrder(c echo.Context) error {
	id := c.Param("id")
	orderRes, err := oh.orderService.GetOrder(c.Request().Context(), &order_service.GetOrderRequest{
		Id: id,
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, orderRes, http.StatusOK)
}

func (oh *OrderHandler) UpdateOrder(c echo.Context) error {
	id := c.Param("id")
	var req order_service.UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}
	req.Id = id
	orderRes, err := oh.orderService.UpdateOrder(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, orderRes, http.StatusOK)
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: order_service.NewOrderService(),
	}
}
