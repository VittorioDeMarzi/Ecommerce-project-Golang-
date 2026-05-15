package products

import (
	"context"

	repo "github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product,(error))
	GetProduct(ctx context.Context, id int64) (repo.Product, (error))
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}


func (s *svc) ListProducts(ctx context.Context) ([]repo.Product,(error)) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) GetProduct(ctx context.Context, id int64) (repo.Product,(error)) {
	return s.repo.GetProduct(ctx, id)
}
