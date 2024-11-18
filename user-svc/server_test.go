package main

import (
	"context"
	"errors"
	"testing"

	usergrpc "github.com/chpter/shared/grpc/user"
	mockdb "github.com/chpter/user-svc/db/mock"
	"github.com/chpter/user-svc/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {

	testCases := []struct {
		name          string
		req           *usergrpc.GetUserByIDRequest
		buildStubs    func(queries *mockdb.MockUserQueries)
		checkResponse func(t *testing.T, res *usergrpc.GetUserByIDResponse, err error)
	}{
		{
			name: "OK",
			req: &usergrpc.GetUserByIDRequest{
				Id: int64(1),
			},
			buildStubs: func(db *mockdb.MockUserQueries) {
				db.EXPECT().
					GetUserById(gomock.Any(), int64(1)).
					Return(&model.User{
						ID:        int64(1),
						FirstName: "Hozain",
						LastName:  "Naim",
						Email:     "hozain@chpter.co",
						Username:  "HozainNaim",
					}, nil)
			},
			checkResponse: func(t *testing.T, res *usergrpc.GetUserByIDResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				user := res.GetData()
				require.Equal(t, "success", res.GetStatus())
				require.Equal(t, "Completed successfully", res.GetMessage())
				require.Equal(t, int64(1), user.Id)
				require.Equal(t, "Hozain", user.FirstName)
				require.Equal(t, "Naim", user.LastName)
				require.Equal(t, "hozain@chpter.co", user.Email)
				require.Equal(t, "HozainNaim", user.Username)
			},
		},

		{
			name: "User Not found",
			req: &usergrpc.GetUserByIDRequest{
				Id: int64(100),
			},
			buildStubs: func(db *mockdb.MockUserQueries) {
				db.EXPECT().
					GetUserById(gomock.Any(), int64(100)).
					Return(nil, errors.New("User not found"))
			},
			checkResponse: func(t *testing.T, res *usergrpc.GetUserByIDResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Equal(t, "failed", res.GetStatus())
				require.Equal(t, "couldn't find user with id: 100", res.GetMessage())
				require.Nil(t, res.GetData())

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			databaseCtrl := gomock.NewController(t)
			defer databaseCtrl.Finish()
			database := mockdb.NewMockUserQueries(databaseCtrl)

			tc.buildStubs(database)
			server := newTestServer(database)

			resp, err := server.GetUserByID(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}
