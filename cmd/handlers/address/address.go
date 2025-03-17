package address_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/address_service"
)

type AddressHandler struct {
	addressService address_service.AddressServiceApi
}

func (a *AddressHandler) Create(c echo.Context) error {
	var req address_service.CreateAddressReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	addressRes, err := a.addressService.CreateAddress(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, addressRes, http.StatusCreated)
}

func (a *AddressHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req address_service.UpdateAddressReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	req.Id = id
	addressRes, err := a.addressService.UpdateAddress(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, addressRes, http.StatusOK)
}

func (a *AddressHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := a.addressService.DeleteAddress(c.Request().Context(), &address_service.DeleteAddressReq{
		Ids: []string{id},
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, "Addresses deleted successfully", http.StatusOK)
}

func (a *AddressHandler) List(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	req := address_service.ListAddressReq{
		Limit:  limitInt,
		Offset: offsetInt,
	}
	listRes, err := a.addressService.ListAddresses(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, listRes, http.StatusOK)
}

func (a *AddressHandler) Get(c echo.Context) error {
	id := c.Param("id")
	addressRes, err := a.addressService.GetAddress(c.Request().Context(), &address_service.GetAddressReq{
		Id: id,
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, addressRes, http.StatusOK)
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		addressService: address_service.NewAddressService(),
	}
}
