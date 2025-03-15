package cart_service

import (
	"context"

	"main.go/common"
	"main.go/common/utils"
	"main.go/internal/models"
	"main.go/internal/repository/cart"
	"main.go/internal/repository/product"
	"main.go/internal/repository/user"
	"main.go/internal/service"
)

// CartServiceApi defines the interface for cart service
type CartServiceApi interface {
	AddToCart(ctx context.Context, req *AddToCartRequest) (*CartRes, *service.ServiceError)
	UpdateCart(ctx context.Context, req *UpdateCartRequest) (*CartRes, *service.ServiceError)
	DeleteCart(ctx context.Context, req *DeleteCartRequest) *service.ServiceError
	ListCart(ctx context.Context, req *ListCartRequest) (ListCartRes, *service.ServiceError)
	GetCart(ctx context.Context, req *GetCartRequest) (*CartRes, *service.ServiceError)
}

// CartService implements CartServiceApi
type CartService struct {
	cartRepo    *cart.CartRepo
	userRepo    *user.UserRepo
	productRepo *product.ProductRepo
}

// **Service Methods**
func (c *CartService) AddToCart(ctx context.Context, req *AddToCartRequest) (*CartRes, *service.ServiceError) {
	currentUser := utils.GetCurrentUser(ctx)
	existingCart, err := c.cartRepo.GetCart(ctx, &cart.GetCartRequest{
		UserId:    currentUser,
		ProductId: req.ProductId,
	})
	if err != nil && err.ErrorCode != common.ErrCodeNotFound {
		return nil, service.HandleRepoErr(err, "Failed to get carts")
	}
	newQuantity := existingCart.Quantity + 1
	if existingCart != nil {
		return c.UpdateCart(ctx, &UpdateCartRequest{
			Quantity: &newQuantity,
			Id:       existingCart.Id,
		})
	}
	newCart := models.NewCart(req.ProductId, req.Quantity, currentUser)
	cart, repoErr := c.cartRepo.AddToCart(ctx, newCart)
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to add to cart")
	}
	return c.getCartResponse(ctx, cart)
}

func (c *CartService) UpdateCart(ctx context.Context, req *UpdateCartRequest) (*CartRes, *service.ServiceError) {
	currentUser := utils.GetCurrentUser(ctx)
	existingCart, repoErr := c.cartRepo.GetCart(ctx, &cart.GetCartRequest{
		ProductId: req.ProductId,
		UserId:    currentUser,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Cart not found")
	}
	existingCart.UpdateQuantity(req.Quantity)
	updatedCart, repoErr := c.cartRepo.UpdateCart(ctx, req.Id, existingCart)
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to update cart")
	}
	return c.getCartResponse(ctx, updatedCart)
}

func (c *CartService) DeleteCart(ctx context.Context, req *DeleteCartRequest) *service.ServiceError {
	repoErr := c.cartRepo.DeleteCart(ctx, req.Ids)
	if repoErr != nil {
		return service.HandleRepoErr(repoErr, "Failed to delete cart")
	}
	return nil
}

func (c *CartService) ListCart(ctx context.Context, req *ListCartRequest) (ListCartRes, *service.ServiceError) {
	currentUser := utils.GetCurrentUser(ctx)
	carts, repoErr := c.cartRepo.ListCart(ctx, &cart.ListCartRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		UserId: currentUser,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to list cart items")
	}
	return c.getListCartResponse(ctx, carts)
}

func (c *CartService) GetCart(ctx context.Context, req *GetCartRequest) (*CartRes, *service.ServiceError) {
	cart, repoErr := c.cartRepo.GetCart(ctx, &cart.GetCartRequest{
		Id: req.Id,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to get cart")
	}
	return c.getCartResponse(ctx, cart)
}

// **Helper function to fetch user data and construct response**
func (c *CartService) getCartResponse(ctx context.Context, cart *models.Cart) (*CartRes, *service.ServiceError) {
	users, err := c.userRepo.List(ctx, &user.ListUserRequest{
		Ids: cart.GetAuditFieldsUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	products, err := c.productRepo.Get(ctx, &product.GetProductRequest{
		Id: cart.ProductId,
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to get product")
	}
	return NewCartRes(cart, users, products), nil
}

// **Helper function to process multiple cart items**
func (c *CartService) getListCartResponse(ctx context.Context, carts models.ListCart) (ListCartRes, *service.ServiceError) {
	users, err := c.userRepo.List(ctx, &user.ListUserRequest{
		Ids: carts.GetUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	products, err := c.productRepo.List(ctx, &product.ListProductRequest{
		Ids: carts.GetProductIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list products")
	}
	return NewListCartRes(carts, users, products), nil
}

// **NewCartService initializes a new CartService instance**
func NewCartService() CartServiceApi {
	return &CartService{
		cartRepo:    cart.NewCartRepo(),
		userRepo:    user.NewUserRepo(),
		productRepo: product.NewProductRepo(),
	}
}
