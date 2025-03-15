package product_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/product_service"
)

type ProductHandler struct {
	productService product_service.ProductServiceApi
}

func (p *ProductHandler) Create(c echo.Context) error {
	var req product_service.CreateProductReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	productRes, err := p.productService.CreateProduct(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, productRes, http.StatusCreated)
}

func (p *ProductHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req product_service.UpdateProductReq
	if err := c.Bind(&req); err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Invalid Request", common.ErrCodeInvalidRequest, http.StatusBadRequest), http.StatusBadRequest)
	}
	req.Id = id
	productRes, err := p.productService.UpdateProduct(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, productRes, http.StatusOK)
}

func (p *ProductHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := p.productService.DeleteProduct(c.Request().Context(), &product_service.DeleteProductReq{
		Ids: []string{id},
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, "Products deleted successfully", http.StatusOK)
}

func (p *ProductHandler) List(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	req := product_service.ListProductReq{
		Limit:  limitInt,
		Offset: offsetInt,
	}
	listRes, err := p.productService.ListProducts(c.Request().Context(), &req)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, listRes, http.StatusOK)
}

func (p *ProductHandler) Get(c echo.Context) error {
	id := c.Param("id")
	productRes, err := p.productService.GetProduct(c.Request().Context(), &product_service.GetProductReq{
		Id: id,
	})
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}
	return common.NewResponse(c, true, productRes, http.StatusOK)
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: product_service.NewProductService(),
	}
}
