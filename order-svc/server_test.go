package main

import (
	"context"
	"errors"
	"testing"

	mockdb "github.com/chpter/order-svc/db/mock"
	"github.com/chpter/order-svc/db/model"
	usermocksvc "github.com/chpter/shared/grpc/mock/user"
	ordergrpc "github.com/chpter/shared/grpc/order"
	usergrpc "github.com/chpter/shared/grpc/user"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateOrderAPI(t *testing.T) {
	newOrderReq := &model.CreateOrderRequest{
		UserID:      int64(14),
		TotalAmount: 197.00,
		Currency:    "USD",
		Items: []*model.OrderItem{
			{
				ID:       int64(101),
				Name:     "Wireless Mouse",
				Quantity: 1,
				Price:    25.50,
				Total:    25.50,
			},
			{
				ID:       int64(102),
				Name:     "Mechanical Keyboard",
				Quantity: 2,
				Price:    85.75,
				Total:    171.50,
			},
		},
	}

	newGrpcOrderRequest := &ordergrpc.CreateOrderRequest{
		UserId:      int64(14),
		TotalAmount: 197.00,
		Currency:    "USD",
		Items: []*ordergrpc.OrderItem{
			{
				Id:       101,
				Name:     "Wireless Mouse",
				Quantity: 1,
				Price:    25.50,
				Total:    25.50,
			},
			{
				Id:       102,
				Name:     "Mechanical Keyboard",
				Quantity: 2,
				Price:    85.75,
				Total:    171.50,
			},
		},
	}

	testCases := []struct {
		name          string
		req           *ordergrpc.CreateOrderRequest
		buildStubs    func(queries *mockdb.MockOrderQueries, userSvc *usermocksvc.MockUserServiceClient)
		checkResponse func(t *testing.T, res *ordergrpc.CreateOrderResponse, err error)
	}{
		{
			name: "OK",
			req:  newGrpcOrderRequest,
			buildStubs: func(database *mockdb.MockOrderQueries, userSvc *usermocksvc.MockUserServiceClient) {
				database.EXPECT().
					CreateOrder(gomock.Any(), newOrderReq).
					Return(&model.CreateOrderResponse{
						OrderID: int64(3),
						UserID:  int64(14),
						Status:  "Shipped",
					}, nil)

				userSvc.EXPECT().GetUserByID(gomock.Any(), &usergrpc.GetUserByIDRequest{
					Id: int64(14),
				}).Return(&usergrpc.GetUserByIDResponse{
					Status:  "success",
					Message: "Completed Successfully",
					Data: &usergrpc.User{
						Id:        int64(14),
						FirstName: "Hozain",
						LastName:  "Naim",
						Email:     "hozain@chpter.co",
						Username:  "Hozain",
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, res *ordergrpc.CreateOrderResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdOrder := res.GetData()
				require.Equal(t, int64(3), createdOrder.OrderId)
				require.Equal(t, newOrderReq.UserID, createdOrder.UserId)
				require.Equal(t, "Shipped", createdOrder.Status)
			},
		},
		{
			name: "User Not Found",
			req: &ordergrpc.CreateOrderRequest{
				UserId:      int64(100),
				TotalAmount: 150.00,
				Currency:    "USD",
				Items: []*ordergrpc.OrderItem{
					{
						Id:       111,
						Name:     "Wireless Mouse",
						Quantity: 2,
						Price:    25.50,
						Total:    50.00,
					},
					{
						Id:       112,
						Name:     "Mechanical Keyboard",
						Quantity: 2,
						Price:    50.00,
						Total:    100.00,
					},
				},
			},
			buildStubs: func(database *mockdb.MockOrderQueries, userSvc *usermocksvc.MockUserServiceClient) {
				userSvc.EXPECT().
					GetUserByID(gomock.Any(), &usergrpc.GetUserByIDRequest{
						Id: int64(100),
					}).
					Return(&usergrpc.GetUserByIDResponse{
						Status:  "failed",
						Message: "couldn't find user with id: 100",
						Data:    nil,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *ordergrpc.CreateOrderResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Nil(t, res.GetData())
				require.Equal(t, "failed", res.GetStatus())
				require.Equal(t, "User not found", res.GetMessage())
			},
		},
		{
			name: "Service unavailable",
			req: &ordergrpc.CreateOrderRequest{
				UserId:      int64(100),
				TotalAmount: 150.00,
				Currency:    "USD",
				Items: []*ordergrpc.OrderItem{
					{
						Id:       111,
						Name:     "Wireless Mouse",
						Quantity: 2,
						Price:    25.50,
						Total:    50.00,
					},
					{
						Id:       112,
						Name:     "Mechanical Keyboard",
						Quantity: 2,
						Price:    50.00,
						Total:    100.00,
					},
				},
			},
			buildStubs: func(database *mockdb.MockOrderQueries, userSvc *usermocksvc.MockUserServiceClient) {
				userSvc.EXPECT().
					GetUserByID(gomock.Any(), &usergrpc.GetUserByIDRequest{
						Id: int64(100),
					}).
					Return(nil, errors.New("service unavailable"))
			},
			checkResponse: func(t *testing.T, res *ordergrpc.CreateOrderResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Nil(t, res.GetData())
				require.Equal(t, "failed", res.GetStatus())
				require.Equal(t, "An error occurred. Please try again!", res.GetMessage())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			databaseCtrl := gomock.NewController(t)
			defer databaseCtrl.Finish()
			database := mockdb.NewMockOrderQueries(databaseCtrl)

			userServiceCtrl := gomock.NewController(t)
			defer userServiceCtrl.Finish()
			userService := usermocksvc.NewMockUserServiceClient(userServiceCtrl)

			tc.buildStubs(database, userService)
			server := newTestServer(database, userService)

			resp, err := server.CreateOrder(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}
