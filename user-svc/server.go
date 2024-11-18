package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/chpter/shared/config"
	usergrpc "github.com/chpter/shared/grpc/user"
	"github.com/chpter/shared/lazyerror"
	"github.com/chpter/user-svc/db"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	usergrpc.UnimplementedUserServiceServer
	database   db.UserQueries
	grpcServer *grpc.Server
	_          interface{} // disallow anonymous interface
}

func startGrpcServer(ctx context.Context, cfg *config.Config) (io.Closer, error) {
	listenAddr := fmt.Sprintf(":%d", cfg.UserServiceServer.Port)

	// setup database
	database, err := db.New(ctx, cfg.UserServiceServer.DSN)
	if err != nil {
		slog.Error("connection to database failed with:", slog.Any("error", err))
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("User service failed to listen with:", slog.Any("error", err))
		return nil, lazyerror.Error(err)
	}

	srv := &server{
		database: database,
	}

	srv.grpcServer = grpc.NewServer()

	usergrpc.RegisterUserServiceServer(srv.grpcServer, srv)

	go func() {
		slog.Info("User Service is listenining on: ", slog.String("addr", listenAddr))
		if err := srv.grpcServer.Serve(listener); err != nil {
			slog.Error("cannot start user service, failed with:  ", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	return srv, nil
}

func (s *server) Close() error {
	slog.Info("shutting down UserService gRPC server...")
	s.grpcServer.GracefulStop()

	return nil
}

func (s *server) GetUserByID(ctx context.Context, payload *usergrpc.GetUserByIDRequest) (*usergrpc.GetUserByIDResponse, error) {
	user, err := s.database.GetUserById(ctx, payload.Id)
	if err != nil {
		slog.Error("couldn't fetch user, failed with: ", slog.Any("error", err))
		return &usergrpc.GetUserByIDResponse{
			Status:  "failed",
			Message: fmt.Sprintf("couldn't find user with id: %d", payload.Id),
			Data:    nil,
		}, nil
	}

	return &usergrpc.GetUserByIDResponse{
		Status:  "Success",
		Message: "Completed successfully",
		Data: &usergrpc.User{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(*user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
		},
	}, nil
}

var (
	_ io.Closer = (*server)(nil)
)
