package cart_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/cart_service"
)

type CartHandler struct {
	cartService cart_service.CartServiceApi
}

func (ch *CartHandler) AddToCart(c echo.Context) error {
	var req cart_service.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	cartRes, err := ch.cartService.AddToCart(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, cartRes, http.StatusCreated)
}

func (ch *CartHandler) UpdateCart(c echo.Context) error {
	id := c.Param("id")
	var req cart_service.UpdateCartRequest
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	req.Id = id
	cartRes, err := ch.cartService.UpdateCart(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, cartRes, http.StatusOK)
}

func (ch *CartHandler) DeleteCart(c echo.Context) error {
	id := c.Param("id")
	err := ch.cartService.DeleteCart(c.Request().Context(), &cart_service.DeleteCartRequest{
		Ids: []string{id},
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, "Cart item deleted successfully", http.StatusOK)
}

func (ch *CartHandler) ListCart(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	req := cart_service.ListCartRequest{
		Limit:  limitInt,
		Offset: offsetInt,
	}
	listRes, err := ch.cartService.ListCart(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, listRes, http.StatusOK)
}

func (ch *CartHandler) GetCart(c echo.Context) error {
	id := c.Param("id")
	cartRes, err := ch.cartService.GetCart(c.Request().Context(), &cart_service.GetCartRequest{
		Id: id,
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, cartRes, http.StatusOK)
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: cart_service.NewCartService(),
	}
}
