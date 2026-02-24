package service

import (
	"context"
	"errors"
	"github.com/go-list-templ/grpc/pkg/paginate"
	"testing"
	"time"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/vo"
	"github.com/go-list-templ/grpc/internal/core/dto"
	"github.com/go-list-templ/grpc/internal/port/mock"
	"github.com/google/uuid"
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

	user := entity.User{
		ID:        vo.UnsafeID(uuid.New()),
		Name:      vo.UnsafeName("test"),
		Email:     vo.UnsafeEmail("example@example.com"),
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users := make([]entity.User, 15)

	for i := 0; i < 15; i++ {
		users[i] = user
	}

	pg := paginate.NewUUIDPaginate("")
	pageToken := pg.GenerateToken(user.ID.Value().String())

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
				NextPageToken: pageToken,
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
				Users:         []dto.User{},
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
			output: dto.UserListOutput{},
			err:    errSome,
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

	tests := []struct {
		name   string
		mock   func()
		input  dto.UserCreateInput
		output dto.UserCreateOutput
		err    error
	}{
		{
			name: "success - create user",
			mock: func() {
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: dto.UserCreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: dto.UserCreateOutput{
				User: dto.UserFromEntity(user),
			},
			err: nil,
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
			input: dto.UserCreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: dto.UserCreateOutput{},
			err:    errSome,
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
			input: dto.UserCreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: dto.UserCreateOutput{},
			err:    errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.Create(context.Background(), tt.input)
			require.ErrorIs(t, err, tt.err)

			require.Equal(t, tt.output.User.Name, got.User.Name)
			require.Equal(t, tt.output.User.Email, got.User.Email)
		})
	}
}
