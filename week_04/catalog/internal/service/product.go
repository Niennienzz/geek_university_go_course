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

func NewProductService(uc *biz.ProductUseCase, logger log.Logger) *ProductService {
	return &ProductService{uc: uc, log: log.NewHelper("service/product", logger)}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.CreateProductReply, error) {
	err := s.uc.Create(ctx, &biz.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	return &v1.CreateProductReply{}, err
}

func (s *ProductService) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductReply, error) {
	p, err := s.uc.Read(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetProductReply{
		Value: &v1.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.UpdateProductReply, error) {
	err := s.uc.Update(ctx, &biz.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	return &v1.UpdateProductReply{}, err
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*v1.DeleteProductReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	return &v1.DeleteProductReply{}, err
}

func (s *ProductService) ListProduct(ctx context.Context, req *v1.ListProductRequest) (*v1.ListProductReply, error) {
	ps, err := s.uc.List(ctx)
	reply := &v1.ListProductReply{}
	for _, p := range ps {
		reply.Values = append(reply.Values, &v1.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return reply, err
}
