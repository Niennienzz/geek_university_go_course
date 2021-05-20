package data

import (
	"catalog/internal/biz"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper("data/greeter", logger),
	}
}

func (r *productRepo) CreateProduct(ctx context.Context, val *biz.Product) error {
	_, err := r.data.db.Product.Create().
		SetName(val.Name).
		SetDescription(val.Description).
		SetPrice(val.Price).
		Save(ctx)
	return err
}

func (r *productRepo) ReadProduct(ctx context.Context, id int64) (*biz.Product, error) {
	p, err := r.data.db.Product.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, val *biz.Product) error {
	p, err := r.data.db.Product.Get(ctx, val.ID)
	if err != nil {
		return err
	}
	_, err = p.Update().
		SetName(val.Name).
		SetDescription(val.Description).
		SetPrice(val.Price).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

func (r *productRepo) DeleteProduct(ctx context.Context, id int64) error {
	return r.data.db.Product.DeleteOneID(id).Exec(ctx)
}

func (r *productRepo) ListArticle(ctx context.Context) ([]*biz.Product, error) {
	ps, err := r.data.db.Product.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Product, 0)
	for _, p := range ps {
		rv = append(rv, &biz.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return rv, nil
}
