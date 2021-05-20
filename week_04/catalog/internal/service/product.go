package service

import (
	"context"

	v1 "catalog/api/catalog/v1"
	"catalog/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ProductService struct {
	v1.UnimplementedCatalogServiceServer

	uc  *biz.ProductUseCase
	log *log.Helper
}

func NewGreeterService(uc *biz.ProductUseCase, logger log.Logger) *ProductService {
	return &ProductService{uc: uc, log: log.NewHelper("service/greeter", logger)}
}

func (s *ProductService) CreateProduct(context.Context, *v1.CreateProductRequest) (*v1.CreateProductReply, error) {
	return nil, nil
}

func (s *ProductService) GetProduct(context.Context, *v1.GetProductRequest) (*v1.GetProductReply, error) {
	return nil, nil
}

func (s *ProductService) UpdateProduct(context.Context, *v1.UpdateProductRequest) (*v1.UpdateProductReply, error) {
	return nil, nil
}

func (s *ProductService) DeleteProduct(context.Context, *v1.DeleteProductRequest) (*v1.DeleteProductReply, error) {
	return nil, nil
}

func (s *ProductService) ListProduct(context.Context, *v1.ListProductRequest) (*v1.ListProductReply, error) {
	return nil, nil
}
