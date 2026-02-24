package service

import (
	"context"
	"errors"
	"github.com/go-list-templ/grpc/internal/core/dto"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/vo"
	"github.com/go-list-templ/grpc/internal/port/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errSome = errors.New("something went wrong")

func mocks(t *testing.T) (*mock.MockUserRepo, *mock.MockOutboxRepo, *mock.MockTransactionManager) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	ur := mock.NewMockUserRepo(mockCtl)
	or := mock.NewMockOutboxRepo(mockCtl)
	tm := mock.NewMockTransactionManager(mockCtl)

	return ur, or, tm
}

func TestNewUser(t *testing.T) {
	t.Run("new user", func(t *testing.T) {
		ur, or, tm := mocks(t)

		got := NewUser(ur, or, tm)
		require.Equal(t, &User{
			ur, or, tm,
		}, got)
	})
}

func TestUser_List(t *testing.T) {
	ur, or, tm := mocks(t)
	userService := NewUser(ur, or, tm)

	users := []entity.User{
		{
			ID:        vo.UnsafeID(uuid.New()),
			Name:      vo.UnsafeName("test"),
			Email:     vo.UnsafeEmail("example@example.com"),
			Avatar:    vo.NewAvatar(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name   string
		mock   func()
		input  dto.UserListInput
		output dto.UserListOutput
		err    error
	}{
		{
			name: "success - get all users",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(users, nil)
			},
			input: dto.UserListInput{
				PageToken: "",
			},
			output: dto.UserListOutput{
				Users:         dto.UsersFromEntity(users),
				NextPageToken: "",
			},
			err: nil,
		},
		{
			name: "success - get empty users",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return([]entity.User{}, nil)
			},
			input: dto.UserListInput{
				PageToken: "",
			},
			output: dto.UserListOutput{
				Users:         dto.UsersFromEntity(users),
				NextPageToken: "",
			},
			err: nil,
		},
		{
			name: "fail - err user repo",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return([]entity.User{}, errSome)
			},
			input: dto.UserListInput{
				PageToken: "",
			},
			output: dto.UserListOutput{
				Users:         dto.UsersFromEntity(users),
				NextPageToken: "",
			},
			err: errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.List(context.Background(), tt.input)

			require.Equal(t, tt.output, got)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestUser_Create(t *testing.T) {
	ur, or, tm := mocks(t)
	userService := NewUser(ur, or, tm)

	user := entity.User{
		ID:        vo.UnsafeID(uuid.New()),
		Name:      vo.UnsafeName("test"),
		Email:     vo.UnsafeEmail("example@example.com"),
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	type args struct {
		user entity.User
	}
	tests := []struct {
		name string
		mock func()
		args args
		want entity.User
		err  error
	}{
		{
			name: "success - create user",
			mock: func() {
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				user: user,
			},
			want: user,
			err:  nil,
		},
		{
			name: "fail - err in user repo",
			mock: func() {
				ur.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errSome)
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
			},
			args: args{
				user: user,
			},
			want: entity.User{},
			err:  errSome,
		},
		{
			name: "fail - err in outbox repo",
			mock: func() {
				ur.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errSome)

				tm.EXPECT().Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
			},
			args: args{
				user: user,
			},
			want: entity.User{},
			err:  errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.Create(context.Background(), tt.args.user)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
