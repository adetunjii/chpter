//go:generate mockgen -source=db.go -destination=./mock/db.go -typed

package db

import (
	"context"
	"errors"

	"github.com/chpter/user-svc/db/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserQueries interface {
	GetUsers(ctx context.Context, limit int64, offset int64) ([]*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
}

type Database struct {
	ctx context.Context
	p   *pgxpool.Pool
	_   interface{} // disallow anonymous fields
}

func New(ctx context.Context, uri string) (*Database, error) {
	return nil, errors.New("not implemented")
}

var (
	_ UserQueries = (*Database)(nil)
)
