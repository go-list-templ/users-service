package handler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// nolint:goconst
func TestCreateUser(t *testing.T) {
	host := "app"
	grpcURL := host + ":8080"

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() {
		err = grpcConn.Close()
		require.NoError(t, err)
	}()

	UserServiceClient := v1.NewUserServiceClient(grpcConn)

	type args struct {
		request *v1.CreateUserRequest
	}
	tests := []struct {
		name string
		args args
		want *v1.CreateUserResponse
		err  error
	}{
		{
			name: "success - create user",
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "test",
					Email: "test@test.com",
				},
			},
			want: &v1.CreateUserResponse{
				User: &v1.User{
					Name:  "test",
					Email: "test@test.com",
				},
			},
			err: nil,
		},
		{
			name: "fail - already exists user",
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "test",
					Email: "test@test.com",
				},
			},
			want: nil,
			err:  status.Error(codes.AlreadyExists, entity.ErrUserAlreadyExists.Error()),
		},
		{
			name: "fail - create invalid name",
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "1",
					Email: "test@test.com",
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, entity.ErrUserAlreadyExists.Error()),
		},
		{
			name: "fail - create invalid email",
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "1",
					Email: "test@",
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, entity.ErrUserAlreadyExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserServiceClient.CreateUser(t.Context(), tt.args.request)
			if tt.want != nil {
				require.NoError(t, err)

				require.NoError(t, uuid.Validate(got.User.Id))
				require.NotEmpty(t, got.User.Avatar)
				require.NotEmpty(t, got.User.CreatedAt)
				require.NotEmpty(t, got.User.UpdatedAt)

				require.Equal(t, tt.want.User.Name, got.User.Name)
				require.Equal(t, tt.want.User.Email, got.User.Email)
			} else {
				require.Error(t, err)
				require.Equal(t, err, tt.err)
			}
		})
	}
}

// nolint:goconst
func TestAllUsers(t *testing.T) {
	host := "app"
	grpcURL := host + ":8080"
	requests := 10

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	defer func() {
		err = grpcConn.Close()
		require.NoError(t, err)
	}()

	UserServiceClient := v1.NewUserServiceClient(grpcConn)

	for i := 0; i < requests; i++ {
		resp, err := UserServiceClient.AllUsers(t.Context(), &v1.AllUsersRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Users, 1)
		require.Equal(t, resp.Users[0].Name, "test")
		require.Equal(t, resp.Users[0].Email, "test@test.com")
	}
}
