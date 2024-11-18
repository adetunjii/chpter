package db

import (
	"context"
	"errors"

	"github.com/chpter/user-svc/db/model"
)

func (d *Database) GetUsers(ctx context.Context, limit int64, offset int64) ([]*model.User, error) {
	return nil, errors.New("not implemented")
}

func (d *Database) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	return nil, errors.New("not implemented")
}
