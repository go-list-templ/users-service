package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port/mock"
	"github.com/go-list-templ/users-service/pkg/paginate"
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

	generateOutput := func(count int, token string) dto.ListOutput {
		users := make([]dto.User, count)

		for i := 0; i < count; i++ {
			users[i] = dto.FromEntity(user)
		}

		return dto.ListOutput{
			Users:         users,
			NextPageToken: token,
		}
	}

	pg := paginate.NewUUIDPaginate("")

	limit := pg.Limit()
	pageToken := pg.GenerateToken(user.ID.Value().String())

	tests := []struct {
		name          string
		mock          func()
		wantPageToken string
		wantCount     int
		err           error
	}{
		{
			name: "success - empty page token len less than 1",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(generateOutput(limit-1, ""), nil)
			},
			wantPageToken: "",
			wantCount:     14,
			err:           nil,
		},
		{
			name: "success - empty page token len equal limit",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(generateOutput(limit, ""), nil)
			},
			wantPageToken: "",
			wantCount:     15,
			err:           nil,
		},
		{
			name: "success - get page token len more by 1",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(generateOutput(limit+1, pageToken), nil)
			},
			wantPageToken: pageToken,
			wantCount:     16,
			err:           nil,
		},
		{
			name: "success - empty result",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(generateOutput(0, ""), nil)
			},
			wantPageToken: "",
			wantCount:     0,
			err:           nil,
		},
		{
			name: "fail - err at get list users",
			mock: func() {
				ur.EXPECT().All(gomock.Any(), gomock.Any()).Return(dto.ListOutput{}, errSome)
			},
			wantPageToken: "",
			wantCount:     0,
			err:           errSome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.List(context.Background(), dto.ListInput{})
			require.ErrorIs(t, err, tt.err)

			require.Equal(t, tt.wantPageToken, got.NextPageToken)
			require.Equal(t, tt.wantCount, len(got.Users))
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
		input  dto.CreateInput
		output entity.User
		isErr  bool
	}{
		{
			name: "success - create user",
			mock: func() {
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: dto.CreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: user,
			isErr:  false,
		},
		{
			name: "fail - min len name",
			mock: func() {},
			input: dto.CreateInput{
				Name:  "t",
				Email: "example@example.com",
			},
			output: user,
			isErr:  true,
		},
		{
			name: "fail - invalid email",
			mock: func() {},
			input: dto.CreateInput{
				Name:  "test",
				Email: "test",
			},
			output: user,
			isErr:  true,
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
			input: dto.CreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: entity.User{},
			isErr:  true,
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
			input: dto.CreateInput{
				Name:  "test",
				Email: "example@example.com",
			},
			output: entity.User{},
			isErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.Create(context.Background(), tt.input)
			if tt.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.output.Name.Value(), got.Name.Value())
				require.Equal(t, tt.output.Email.Value(), got.Email.Value())
			}
		})
	}
}
