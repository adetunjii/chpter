package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/chpter/order-svc/db"
	"github.com/chpter/order-svc/db/model"
	"github.com/chpter/shared/config"
	ordergrpc "github.com/chpter/shared/grpc/order"
	usergrpc "github.com/chpter/shared/grpc/user"
	"github.com/chpter/shared/lazyerror"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	ordergrpc.UnimplementedOrderServiceServer
	userSvc    usergrpc.UserServiceClient
	database   db.OrderQueries
	grpcServer *grpc.Server
	_          interface{} // disallow anonymous interface
}

func startGrpcServer(ctx context.Context, cfg *config.Config) (io.Closer, error) {
	listenAddr := fmt.Sprintf(":%d", cfg.OrderServiceServer.Port)

	// setup database
	database, err := db.New(ctx, cfg.OrderServiceServer.DSN)
	if err != nil {
		slog.Error("connection to database failed with:", slog.Any("error", err))
		os.Exit(1)
	}

	userServiceURL := fmt.Sprintf("%s:%d", cfg.UserServiceServer.Host, cfg.UserServiceServer.Port)
	clientConn, err := grpc.Dial(userServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("connection to user service failed with:", slog.Any("error", err))
		return nil, lazyerror.Error(err)
	}
	defer clientConn.Close()

	userClient := usergrpc.NewUserServiceClient(clientConn)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("User service failed to listen with:", slog.Any("error", err))
		return nil, lazyerror.Error(err)
	}

	srv := &server{
		userSvc:  userClient,
		database: database,
	}

	srv.grpcServer = grpc.NewServer()

	ordergrpc.RegisterOrderServiceServer(srv.grpcServer, srv)

	go func() {
		slog.Info("Order Service is listenining on: ", slog.String("addr", listenAddr))
		if err := srv.grpcServer.Serve(listener); err != nil {
			slog.Error("cannot start order service, failed with:  ", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	return srv, nil
}

func (s *server) Close() error {
	slog.Info("shutting down OrderService gRPC server...")

	s.grpcServer.GracefulStop()

	return nil
}

func mapOrderItems(items []*ordergrpc.OrderItem) []*model.OrderItem {
	orderItems := make([]*model.OrderItem, len(items))

	for i := 0; i < len(items); i++ {
		orderItems[i] = &model.OrderItem{
			ID:       items[i].Id,
			Name:     items[i].Name,
			Quantity: int(items[i].Quantity),
			Price:    items[i].Price,
			Total:    items[i].Total,
		}
	}

	return orderItems
}

func mapToGrpcOrderItems(items []*model.OrderItem) []*ordergrpc.OrderItem {
	orderItems := make([]*ordergrpc.OrderItem, len(items))

	for i := 0; i < len(items); i++ {
		orderItems[i] = &ordergrpc.OrderItem{
			Id:       items[i].ID,
			Name:     items[i].Name,
			Quantity: int32(items[i].Quantity),
			Price:    items[i].Price,
			Total:    items[i].Total,
		}
	}

	return orderItems
}

func mapToGrpcOrders(items []*model.Order) []*ordergrpc.Order {
	orders := make([]*ordergrpc.Order, len(items))

	for i := 0; i < len(items); i++ {
		orders[i] = &ordergrpc.Order{
			Id:          items[i].ID,
			TotalAmount: items[i].TotalAmount,
			UserId:      items[i].UserID,
			Currency:    items[i].Currency,
			Status:      items[i].Status,
			Items:       mapToGrpcOrderItems(items[i].Items),
			CreatedAt:   timestamppb.New(*items[i].CreatedAt),
			UpdatedAt:   timestamppb.New(*items[i].UpdatedAt),
		}
	}

	return orders
}

func (s *server) CreateOrder(ctx context.Context, reqPayload *ordergrpc.CreateOrderRequest) (*ordergrpc.CreateOrderResponse, error) {
	var wg sync.WaitGroup
	wg.Add(1)

	payload := &model.CreateOrderRequest{
		UserID:      reqPayload.UserId,
		TotalAmount: reqPayload.TotalAmount,
		Currency:    reqPayload.Currency,
	}

	payload.Items = mapOrderItems(reqPayload.Items)

	var userResp *usergrpc.GetUserByIDResponse
	var userErr error

	// fetch user concurrently from UserService
	go func() {
		defer wg.Done()
		userResp, userErr = s.userSvc.GetUserByID(ctx, &usergrpc.GetUserByIDRequest{
			Id: reqPayload.UserId,
		})
	}()

	wg.Wait()

	if userErr != nil {
		slog.Error("User service failed with:", slog.Any("error", userErr))
		return &ordergrpc.CreateOrderResponse{
			Status:  "failed",
			Message: "An error occurred. Please try again!",
			Data:    nil,
		}, nil
	}

	if userResp.GetStatus() == "failed" {
		slog.Error("couldn't fetch user info, failed with: ", slog.Any("error", userResp.GetMessage()))
		return &ordergrpc.CreateOrderResponse{
			Status:  "failed",
			Message: "User not found",
			Data:    nil,
		}, nil
	}

	slog.Info("Found", slog.Any("user", userResp))

	createdOrder, err := s.database.CreateOrder(ctx, payload)
	if err != nil {
		slog.Error("cannot create order, failed with: ", slog.Any("error", err))
		return &ordergrpc.CreateOrderResponse{
			Status:  "failed",
			Message: "Order creation did not complete successfully",
			Data:    nil,
		}, nil
	}

	return &ordergrpc.CreateOrderResponse{
		Status:  "success",
		Message: "Completed successfully",
		Data: &ordergrpc.CreatedOrder{
			OrderId: createdOrder.OrderID,
			UserId:  userResp.Data.Id,
			Status:  "Shipped",
		},
	}, nil
}

func (s *server) GetOrderByUserID(ctx context.Context, reqPayload *ordergrpc.GetOrderByUserIDRequest) (*ordergrpc.GetOrdersByUserIDResponse, error) {
	orders, err := s.database.GetOrdersByUserID(ctx, reqPayload.UserId)
	if err != nil {
		slog.Error("cannot fetch user' orders, failed with: ", slog.Any("error", err))
		return &ordergrpc.GetOrdersByUserIDResponse{
			Status:  "failed",
			Message: "couldn't fetch users' orders",
			Data:    nil,
		}, nil
	}

	return &ordergrpc.GetOrdersByUserIDResponse{
		Status:  "success",
		Message: "Completed successfully",
		Data:    mapToGrpcOrders(orders),
	}, nil
}

var (
	_ io.Closer = (*server)(nil)
)
