package db

import (
	"context"
	"errors"

	"github.com/chpter/order-svc/db/model"
)

func (d *Database) GetOrders(ctx context.Context, limit int64, offset int64) ([]*model.Order, error) {
	return nil, errors.New("not implemented")
}

func (d *Database) GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	return nil, errors.New("not implemented")
}

func (d *Database) GetOrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	return nil, errors.New("not implemented")
}

func (d *Database) CreateOrder(ctx context.Context, payload *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	return nil, errors.New("not implemented")
}
