package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error)
	GetProductsByIds(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{repository: r}
}

func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	p := Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
	}

	err := s.repository.PutProduct(ctx, p)
	return &p, err
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListProduct(ctx, skip, take)
}

func (s *catalogService) GetProductsByIds(ctx context.Context, ids []string) ([]*Product, error) {
	return s.repository.ListProductsWithIDS(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.SearchProduct(ctx, query, skip, take)
}
