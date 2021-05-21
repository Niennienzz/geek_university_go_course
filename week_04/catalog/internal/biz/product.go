package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductRepo interface {
	CreateProduct(context.Context, *Product) error
	ReadProduct(ctx context.Context, id int64) (*Product, error)
	UpdateProduct(context.Context, *Product) error
	DeleteProduct(ctx context.Context, id int64) error
	ListArticle(ctx context.Context) ([]*Product, error)
}

type ProductUseCase struct {
	repo ProductRepo
	log  *log.Helper
}

func NewProductUseCase(repo ProductRepo, logger log.Logger) *ProductUseCase {
	return &ProductUseCase{repo: repo, log: log.NewHelper("usecase/product", logger)}
}

func (uc *ProductUseCase) Create(ctx context.Context, val *Product) error {
	return uc.repo.CreateProduct(ctx, val)
}

func (uc *ProductUseCase) Read(ctx context.Context, id int64) (*Product, error) {
	return uc.repo.ReadProduct(ctx, id)
}

func (uc *ProductUseCase) Update(ctx context.Context, val *Product) error {
	return uc.repo.UpdateProduct(ctx, val)
}

func (uc *ProductUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteProduct(ctx, id)
}

func (uc *ProductUseCase) List(ctx context.Context) ([]*Product, error) {
	return uc.repo.ListArticle(ctx)
}
