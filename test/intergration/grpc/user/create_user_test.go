package user

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateUser(t *testing.T) {
	host := "app"
	grpcURL := host + ":8080"
	requestTimeout := 1 * time.Second

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	defer func() {
		err = grpcConn.Close()
		require.NoError(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

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
			name: "fail - create invalid user",
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "",
					Email: "",
				},
			},
			want: nil,
			err:  errors.New("user name is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserServiceClient.CreateUser(ctx, tt.args.request)
			require.Equal(t, tt.err, err)

			require.NoError(t, uuid.Validate(got.User.Id))

			require.NotEmpty(t, got.User.Avatar)
			require.NotEmpty(t, got.User.CreatedAt)
			require.NotEmpty(t, got.User.UpdatedAt)

			require.Equal(t, tt.want.User.Name, got.User.Name)
			require.Equal(t, tt.want.User.Email, got.User.Email)
		})
	}
}
