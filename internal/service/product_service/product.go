package product_service

import (
	"context"

	"main.go/internal/models"
	"main.go/internal/repository/product"
	"main.go/internal/repository/user"
	"main.go/internal/service"
)

type ProductServiceApi interface {
	CreateProduct(ctx context.Context, req *CreateProductReq) (*ProductRes, *service.ServiceError)
	UpdateProduct(ctx context.Context, req *UpdateProductReq) (*ProductRes, *service.ServiceError)
	DeleteProduct(ctx context.Context, req *DeleteProductReq) *service.ServiceError
	GetProduct(ctx context.Context, req *GetProductReq) (*ProductRes, *service.ServiceError)
	ListProducts(ctx context.Context, req *ListProductReq) (ListProductRes, *service.ServiceError)
}

type ProductService struct {
	productRepo *product.ProductRepo
	userRepo    *user.UserRepo
}

// CreateProduct implements ProductServiceApi.
func (p *ProductService) CreateProduct(ctx context.Context, req *CreateProductReq) (*ProductRes, *service.ServiceError) {
	product, repoErr := p.productRepo.CreateProduct(ctx, models.NewProduct(req.Name, req.Images, req.Price))
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to create product")
	}
	return p.getProductResponse(ctx, product)
}

func (p *ProductService) getProductResponse(ctx context.Context, product *models.Product) (*ProductRes, *service.ServiceError) {
	users, err := p.userRepo.List(ctx, &user.ListUserRequest{
		Ids: product.GetAuditFieldsUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	return NewProductRes(product, users), nil
}

// UpdateProduct implements ProductServiceApi.
func (p *ProductService) UpdateProduct(ctx context.Context, req *UpdateProductReq) (*ProductRes, *service.ServiceError) {
	existingProduct, repoErr := p.productRepo.Get(ctx, &product.GetProductRequest{Id: req.Id})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Product not found")
	}
	existingProduct.UpdateName(req.Name)
	existingProduct.UpdatePrice(req.Price)
	existingProduct.UpdateImages(req.Images)

	updatedProduct, repoErr := p.productRepo.Save(ctx, existingProduct)
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to update product")
	}
	return p.getProductResponse(ctx, updatedProduct)
}

// DeleteProduct implements ProductServiceApi.
func (p *ProductService) DeleteProduct(ctx context.Context, req *DeleteProductReq) *service.ServiceError {
	repoErr := p.productRepo.Delete(ctx, req.Ids)
	if repoErr != nil {
		return service.HandleRepoErr(repoErr, "Failed to delete product")
	}
	return nil
}

// GetProduct implements ProductServiceApi.
func (p *ProductService) GetProduct(ctx context.Context, req *GetProductReq) (*ProductRes, *service.ServiceError) {
	product, repoErr := p.productRepo.Get(ctx, &product.GetProductRequest{Id: req.Id, Name: req.Name})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to get product")
	}
	return p.getProductResponse(ctx, product)
}

// ListProducts implements ProductServiceApi.
func (p *ProductService) ListProducts(ctx context.Context, req *ListProductReq) (ListProductRes, *service.ServiceError) {
	products, repoErr := p.productRepo.List(ctx, &product.ListProductRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Ids:    req.Ids,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to list products")
	}
	return p.getListProductres(ctx, products)
}

func (p *ProductService) getListProductres(ctx context.Context, products models.ListProductResponse) (ListProductRes, *service.ServiceError) {
	users, err := p.userRepo.List(ctx, &user.ListUserRequest{
		Ids: products.GetUserIds(),
	})
	if err != nil {
		return nil, service.HandleRepoErr(err, "Failed to list users")
	}
	return NewListProductRes(products, users), nil
}

// NewProductService initializes a new ProductService instance
func NewProductService() ProductServiceApi {
	return &ProductService{
		productRepo: product.NewProductRepo(),
		userRepo:    user.NewUserRepo(),
	}
}
