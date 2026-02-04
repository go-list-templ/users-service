package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/vo"
	"github.com/go-list-templ/grpc/internal/port/mock"
	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errSome = errors.New("something went wrong")

func mocks(t *testing.T) (*mock.MockUserService, *zap.Logger) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mock.NewMockUserService(ctrl)
	l := zap.NewNop()

	return us, l
}

func TestUser_CreateUser(t *testing.T) {
	us, l := mocks(t)

	userHandler := &User{
		userService: us,
		logger:      l,
	}

	domUser := entity.User{
		ID:        vo.NewID(),
		Name:      vo.UnsafeName("John"),
		Email:     vo.UnsafeEmail("john@example.com"),
		Avatar:    vo.UnsafeAvatar(""),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	type args struct {
		request *v1.CreateUserRequest
	}
	tests := []struct {
		name string
		mock func()
		args args
		want *v1.CreateUserResponse
		err  error
	}{
		{
			name: "success - create user",
			mock: func() {
				us.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domUser, nil)
			},
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "John",
					Email: "john@example.com",
				},
			},
			want: &v1.CreateUserResponse{
				User: &v1.User{
					Id:        domUser.ID.Value().String(),
					Name:      domUser.Name.Value(),
					Email:     domUser.Email.Value(),
					Avatar:    domUser.Avatar.Value(),
					CreatedAt: timestamppb.New(domUser.CreatedAt),
					UpdatedAt: timestamppb.New(domUser.UpdatedAt),
				},
			},
			err: nil,
		},
		{
			name: "fail - invalid request",
			mock: func() {},
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "",
					Email: "",
				},
			},
			want: nil,
			err:  vo.ErrNameMinLength,
		},
		{
			name: "fail - err at create user",
			mock: func() {
				us.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entity.User{}, errSome)
			},
			args: args{
				request: &v1.CreateUserRequest{
					Name:  "John",
					Email: "john@example.com",
				},
			},
			want: nil,
			err:  errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userHandler.CreateUser(context.Background(), tt.args.request)

			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUser_AllUsers(t *testing.T) {
	us, l := mocks(t)

	userHandler := &User{
		userService: us,
		logger:      l,
	}

	domUsers := []entity.User{
		{
			ID:        vo.NewID(),
			Name:      vo.UnsafeName("John"),
			Email:     vo.UnsafeEmail("john@example.com"),
			Avatar:    vo.UnsafeAvatar(""),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        vo.NewID(),
			Name:      vo.UnsafeName("John2"),
			Email:     vo.UnsafeEmail("john2@example.com"),
			Avatar:    vo.UnsafeAvatar(""),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	tests := []struct {
		name string
		mock func()
		want *v1.AllUsersResponse
		err  error
	}{
		{
			name: "success - get all users",
			mock: func() {
				us.EXPECT().All(gomock.Any()).Return(domUsers, nil)
			},
			want: &v1.AllUsersResponse{
				Users: []*v1.User{
					{
						Id:        domUsers[0].ID.Value().String(),
						Name:      domUsers[0].Name.Value(),
						Email:     domUsers[0].Email.Value(),
						Avatar:    domUsers[0].Avatar.Value(),
						CreatedAt: timestamppb.New(domUsers[0].CreatedAt),
						UpdatedAt: timestamppb.New(domUsers[0].UpdatedAt),
					},
					{
						Id:        domUsers[1].ID.Value().String(),
						Name:      domUsers[1].Name.Value(),
						Email:     domUsers[1].Email.Value(),
						Avatar:    domUsers[1].Avatar.Value(),
						CreatedAt: timestamppb.New(domUsers[1].CreatedAt),
						UpdatedAt: timestamppb.New(domUsers[1].UpdatedAt),
					},
				},
			},
			err: nil,
		},
		{
			name: "fail - err at all user",
			mock: func() {
				us.EXPECT().All(gomock.Any()).Return([]entity.User{}, errSome)
			},
			want: nil,
			err:  errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userHandler.AllUsers(context.Background(), &v1.AllUsersRequest{})

			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUser_toProto(t *testing.T) {
	us, l := mocks(t)

	userHandler := &User{
		userService: us,
		logger:      l,
	}

	domUser := entity.User{
		ID:        vo.NewID(),
		Name:      vo.UnsafeName("John"),
		Email:     vo.UnsafeEmail("john@example.com"),
		Avatar:    vo.UnsafeAvatar(""),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	emptyUser := entity.User{}

	type args struct {
		user entity.User
	}
	tests := []struct {
		name string
		args args
		want *v1.User
	}{
		{
			name: "success - to proto entity",
			args: args{
				user: domUser,
			},
			want: &v1.User{
				Id:        domUser.ID.Value().String(),
				Name:      domUser.Name.Value(),
				Email:     domUser.Email.Value(),
				Avatar:    domUser.Avatar.Value(),
				CreatedAt: timestamppb.New(domUser.CreatedAt),
				UpdatedAt: timestamppb.New(domUser.UpdatedAt),
			},
		},
		{
			name: "success - to proto empty entity",
			args: args{
				user: emptyUser,
			},
			want: &v1.User{
				Id:        emptyUser.ID.Value().String(),
				Name:      emptyUser.Name.Value(),
				Email:     emptyUser.Email.Value(),
				Avatar:    emptyUser.Avatar.Value(),
				CreatedAt: timestamppb.New(emptyUser.CreatedAt),
				UpdatedAt: timestamppb.New(emptyUser.UpdatedAt),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userHandler.toProto(tt.args.user)
			require.Equal(t, tt.want, got)
		})
	}
}
