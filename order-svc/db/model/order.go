package model

import "time"

type OrderItem struct {
	ID       int64
	Name     string
	Quantity int
	Price    float64
	Total    float64
}

type Order struct {
	ID          int64        `json:"id"`
	TotalAmount float64      `json:"total_amount"`
	UserID      int64        `json:"user_id"`
	Currency    string       `json:"currency"`
	Status      string       `json:"status"`
	Items       []*OrderItem `json:"items"`
	CreatedAt   *time.Time   `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID      int64        `json:"user_id"`
	TotalAmount float64      `json:"total_amount"`
	Currency    string       `json:"currency"`
	Items       []*OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Status  string `json:"status"`
}
