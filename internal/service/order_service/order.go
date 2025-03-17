package order_service

import (
	"context"

	"main.go/common/utils"
	"main.go/internal/models"
	"main.go/internal/repository/order"
	"main.go/internal/repository/product"
	"main.go/internal/service"
)

type OrderServiceApi interface {
	AddOrder(ctx context.Context, req *AddOrderRequest) *service.ServiceError
	ListOrders(ctx context.Context, req *ListOrderRequest) ([]*OrderResponse, *service.ServiceError)
	GetOrder(ctx context.Context, req *GetOrderRequest) (*OrderResponse, *service.ServiceError)
	UpdateOrder(ctx context.Context, req *UpdateOrderRequest) (*OrderResponse, *service.ServiceError)
}

type OrderService struct {
	orderRepo        *order.OrderRepo
	orderProductRepo *order.OrderProductRepo
	productRepo      *product.ProductRepo
}

var orderService *OrderService

func (o *OrderService) AddOrder(ctx context.Context, req *AddOrderRequest) *service.ServiceError {
	userId := utils.GetCurrentUser(ctx)
	productIds := req.GetProductIds()
	newOrder := models.NewOrder(userId, 0)
	products, err := o.productRepo.List(ctx, &product.ListProductRequest{
		Ids: productIds,
	})
	if err != nil {
		return service.HandleRepoErr(err, "Failed to list products")
	}
	productList := models.ListProductResponse(products)
	productMap := productList.ToMap()
	newOrderProducts := []*models.OrderProduct{}
	createOrder, err := o.orderRepo.CreateOrder(ctx, newOrder)
	if err != nil {
		return service.HandleRepoErr(err, "Failed to create order")
	}
	for _, prod := range req.Products {
		product := productMap[prod.Id]
		if product != nil {
			newOrderProducts = append(newOrderProducts,
				models.NewOrderProductWithOrderId(prod.Id, prod.Quantity, product.Price, createOrder.Id))
		}
	}
	if err := o.orderProductRepo.BulkCreate(ctx, newOrderProducts); err != nil {
		return service.HandleRepoErr(err, "Failed to create order products")
	}
	return nil
}

func (o *OrderService) ListOrders(ctx context.Context, req *ListOrderRequest) ([]*OrderResponse, *service.ServiceError) {
	userId := utils.GetCurrentUser(ctx)
	orders, err := o.orderRepo.List(ctx, &order.ListOrderRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		UserId: userId,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list orders")
	}
	orderList := models.ListOrderResponse(orders)
	orderIds := orderList.GetIds()
	orderProducts, err := o.orderProductRepo.List(ctx, &order.ListOrderProductRequest{
		OrderIds: orderIds,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list order products")
	}
	orderProductList := models.ListOrderProductResponse(orderProducts)
	products, err := o.productRepo.List(ctx, &product.ListProductRequest{
		Ids: orderProductList.ListProductIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list products")
	}
	productList := models.ListProductResponse(products)
	return NewOrderListResponse(orders, productList.ToMap(), orderProductList.ToOrderMap()), nil
}

func (o *OrderService) GetOrder(ctx context.Context, req *GetOrderRequest) (*OrderResponse, *service.ServiceError) {
	orderObj, err := o.orderRepo.Get(ctx, &order.GetOrderRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to get order")
	}
	orderProducts, err := o.orderProductRepo.List(ctx, &order.ListOrderProductRequest{
		OrderId: orderObj.Id,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list order products")
	}
	orderProdList := models.ListOrderProductResponse(orderProducts)
	productIds := orderProdList.ListProductIds()
	products, err := o.productRepo.List(ctx, &product.ListProductRequest{
		Ids: productIds,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list products")
	}
	productList := models.ListProductResponse(products)
	productMap := productList.ToMap()
	return NewOrderResponse(orderObj, productMap, orderProducts), nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, req *UpdateOrderRequest) (*OrderResponse, *service.ServiceError) {
	orderObj, err := o.orderRepo.Get(ctx, &order.GetOrderRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to get order")
	}
	orderObj.UpdateStatus(req.Status)
	_, err = o.orderRepo.Save(ctx, orderObj)
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to update order")
	}
	return o.GetOrder(ctx, &GetOrderRequest{
		Id: req.Id,
	})
}

func NewOrderService() OrderServiceApi {
	if orderService == nil {
		orderService = &OrderService{
			orderRepo:        order.NewOrderRepo(),
			orderProductRepo: order.NewOrderProductRepo(),
			productRepo:      product.NewProductRepo(),
		}
	}
	return orderService
}
