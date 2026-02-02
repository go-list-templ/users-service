package user

import (
	"testing"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
		_, err = UserServiceClient.AllUsers(t.Context(), &v1.AllUsersRequest{})
		require.NoError(t, err)
	}
}
