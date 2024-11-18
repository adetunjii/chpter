package main

import (
	"github.com/chpter/order-svc/db"
	usergrpc "github.com/chpter/shared/grpc/user"
)

func newTestServer(database db.OrderQueries, userSvc usergrpc.UserServiceClient) *server {
	return &server{
		userSvc:  userSvc,
		database: database,
	}
}
