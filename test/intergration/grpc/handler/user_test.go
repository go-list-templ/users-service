package handler

import (
	"strconv"
	"testing"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"

	"github.com/google/uuid"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
		request *v1.CreateRequest
	}
	tests := []struct {
		name string
		args args
		want *v1.CreateResponse
		err  codes.Code
	}{
		{
			name: "success - create user",
			args: args{
				request: &v1.CreateRequest{
					Name:     mo.Some("test").ToPointer(),
					Email:    "test@test.com",
					Password: "password",
				},
			},
			want: &v1.CreateResponse{
				User: &v1.User{
					Name:  mo.Some("test").ToPointer(),
					Email: "test@test.com",
				},
			},
			err: codes.OK,
		},
		{
			name: "success - create with empty name",
			args: args{
				request: &v1.CreateRequest{
					Name:     nil,
					Email:    "new@test.com",
					Password: "password",
				},
			},
			want: &v1.CreateResponse{
				User: &v1.User{
					Name:  nil,
					Email: "new@test.com",
				},
			},
			err: codes.OK,
		},
		{
			name: "fail - already exists user",
			args: args{
				request: &v1.CreateRequest{
					Name:     mo.Some("test").ToPointer(),
					Email:    "test@test.com",
					Password: "password",
				},
			},
			want: nil,
			err:  codes.AlreadyExists,
		},
		{
			name: "fail - create invalid name",
			args: args{
				request: &v1.CreateRequest{
					Name:     mo.Some("t").ToPointer(),
					Email:    "test@test.com",
					Password: "password",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
		{
			name: "fail - create invalid email",
			args: args{
				request: &v1.CreateRequest{
					Name:     mo.Some("test").ToPointer(),
					Email:    "test",
					Password: "password",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
		{
			name: "fail - create invalid password",
			args: args{
				request: &v1.CreateRequest{
					Name:     mo.Some("test").ToPointer(),
					Email:    "test@test.com",
					Password: "pass",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserServiceClient.Create(t.Context(), tt.args.request)
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

				st := status.Convert(err)
				require.Equal(t, st.Code(), tt.err)
			}
		})
	}
}

// nolint:goconst
func TestListUsers(t *testing.T) {
	host := "app"
	grpcURL := host + ":8080"

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	defer func() {
		err = grpcConn.Close()
		require.NoError(t, err)
	}()

	UserServiceClient := v1.NewUserServiceClient(grpcConn)

	for i := 0; i < 30; i++ {
		_, err := UserServiceClient.Create(t.Context(), &v1.CreateRequest{
			Name:     nil,
			Email:    "test" + strconv.Itoa(i) + "@test.com",
			Password: "password",
		})

		require.NoError(t, err)
	}

	pageToken := ""

	for {
		req := &v1.ListRequest{PageToken: pageToken}
		resp, err := UserServiceClient.List(t.Context(), req)
		require.NoError(t, err)

		if resp.NextPageToken == "" {
			break
		}

		pageToken = resp.NextPageToken
	}
}

// nolint:goconst
func TestGetByEmail(t *testing.T) {
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
		request *v1.GetByEmailRequest
	}
	tests := []struct {
		name string
		args args
		want *v1.GetByEmailResponse
		err  codes.Code
	}{
		{
			name: "success - get user by email",
			args: args{
				request: &v1.GetByEmailRequest{
					Email: "test@test.com",
				},
			},
			want: &v1.GetByEmailResponse{
				User: &v1.User{
					Email: "test@test.com",
				},
			},
			err: codes.OK,
		},
		{
			name: "fail - not found user",
			args: args{
				request: &v1.GetByEmailRequest{
					Email: "notfound@test.com",
				},
			},
			want: nil,
			err:  codes.NotFound,
		},
		{
			name: "fail - invalid email",
			args: args{
				request: &v1.GetByEmailRequest{
					Email: "test",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserServiceClient.GetByEmail(t.Context(), tt.args.request)
			if tt.want != nil {
				require.NoError(t, err)

				require.Equal(t, tt.want.User.Email, got.User.Email)
			} else {
				require.Error(t, err)

				st := status.Convert(err)
				require.Equal(t, st.Code(), tt.err)
			}
		})
	}
}

// nolint:goconst
func TestVerifyCred(t *testing.T) {
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
		request *v1.VerifyCredRequest
	}
	tests := []struct {
		name string
		args args
		want *v1.VerifyCredResponse
		err  codes.Code
	}{
		{
			name: "success - verify credential is success",
			args: args{
				request: &v1.VerifyCredRequest{
					Email:    "test@test.com",
					Password: "password",
				},
			},
			want: &v1.VerifyCredResponse{
				UserId: "",
			},
			err: codes.OK,
		},
		{
			name: "fail - not found user",
			args: args{
				request: &v1.VerifyCredRequest{
					Email:    "notfound@test.com",
					Password: "password",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
		{
			name: "fail - wrong password",
			args: args{
				request: &v1.VerifyCredRequest{
					Email:    "test@test.com",
					Password: "wrong-password",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
		{
			name: "fail - invalid email",
			args: args{
				request: &v1.VerifyCredRequest{
					Email:    "test",
					Password: "password",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
		{
			name: "fail - invalid password",
			args: args{
				request: &v1.VerifyCredRequest{
					Email:    "test@test.com",
					Password: "pass",
				},
			},
			want: nil,
			err:  codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserServiceClient.VerifyCred(t.Context(), tt.args.request)
			if tt.want != nil {
				require.NoError(t, err)

				require.NotEmpty(t, got.UserId)
			} else {
				require.Error(t, err)

				st := status.Convert(err)
				require.Equal(t, st.Code(), tt.err)
			}
		})
	}
}
