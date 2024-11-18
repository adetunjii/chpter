//go:generate mockgen -source=db.go -destination=./mock/db.go -typed

package db

import (
	"context"

	"github.com/chpter/order-svc/db/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQueries interface {
	GetOrders(ctx context.Context, limit int64, offset int64) ([]*model.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error)
	CreateOrder(ctx context.Context, payload *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
}

type Database struct {
	ctx context.Context
	p   *pgxpool.Pool
	_   interface{} // disallow anonymous fields
}

func New(ctx context.Context, uri string) (*Database, error) {
	// _, err := openDB(uri)
	// if err != nil {
	// 	return nil, lazyerror.Error(err)
	// }

	// return &Database{
	// 	ctx: ctx,
	// 	p:   pool,
	// }, nil

	return nil, nil
}

var (
	_ OrderQueries = (*Database)(nil)
)
