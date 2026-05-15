package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorProductNotFound = errors.New("Product not found")
	ErrorProductOutOfStock = errors.New("Product out of stock")
)

type Service interface {
	CreateOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}


type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) CreateOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	// validate order
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer Id is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		// look for the product if exists
		product, err := qtx.GetProduct(ctx, item.ProductId)
		if err != nil {
			return repo.Order{}, ErrorProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrorProductOutOfStock
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    product.ID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
