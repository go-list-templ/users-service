package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port/mock"
	"github.com/go-list-templ/users-service/pkg/paginate"
	"github.com/google/uuid"
	"github.com/samber/mo"
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
		Name:      mo.Some(vo.UnsafeName("test")),
		Password:  vo.UnsafePasswordHash("hash"),
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
				ur.EXPECT().List(gomock.Any(), gomock.Any()).Return(generateOutput(limit-1, ""), nil)
			},
			wantPageToken: "",
			wantCount:     14,
			err:           nil,
		},
		{
			name: "success - empty page token len equal limit",
			mock: func() {
				ur.EXPECT().List(gomock.Any(), gomock.Any()).Return(generateOutput(limit, ""), nil)
			},
			wantPageToken: "",
			wantCount:     15,
			err:           nil,
		},
		{
			name: "success - get page token len more by 1",
			mock: func() {
				ur.EXPECT().List(gomock.Any(), gomock.Any()).Return(generateOutput(limit+1, pageToken), nil)
			},
			wantPageToken: pageToken,
			wantCount:     16,
			err:           nil,
		},
		{
			name: "success - empty result",
			mock: func() {
				ur.EXPECT().List(gomock.Any(), gomock.Any()).Return(generateOutput(0, ""), nil)
			},
			wantPageToken: "",
			wantCount:     0,
			err:           nil,
		},
		{
			name: "fail - err at get list users",
			mock: func() {
				ur.EXPECT().List(gomock.Any(), gomock.Any()).Return(dto.ListOutput{}, errSome)
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

	tests := []struct {
		name    string
		args    dto.CreateInput
		mock    func()
		wantErr bool
	}{
		{
			name: "success - create user",
			mock: func() {
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: dto.CreateInput{
				Name:     mo.Some("test").ToPointer(),
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "success - create user without name",
			mock: func() {
				tm.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: dto.CreateInput{
				Name:     nil,
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "fail - invalid name",
			mock: func() {},
			args: dto.CreateInput{
				Name:     mo.Some("t").ToPointer(),
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "fail - invalid password",
			mock: func() {},
			args: dto.CreateInput{
				Name:     mo.Some("t").ToPointer(),
				Email:    "example@example.com",
				Password: "test",
			},
			wantErr: true,
		},
		{
			name: "fail - invalid email",
			mock: func() {},
			args: dto.CreateInput{
				Name:     mo.Some("test").ToPointer(),
				Email:    "test",
				Password: "password",
			},
			wantErr: true,
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
			args: dto.CreateInput{
				Name:     mo.Some("test").ToPointer(),
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: true,
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
			args: dto.CreateInput{
				Name:     mo.Some("test").ToPointer(),
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.Create(context.Background(), tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var username *string
				if name, ok := got.Name.Get(); ok {
					str := name.Value()
					username = &str
				}

				require.Equal(t, tt.args.Name, username)
				require.Equal(t, tt.args.Email, got.Email.Value())
				require.True(t, got.Password.Compare(vo.UnsafePlainPassword(tt.args.Password)))
			}
		})
	}
}

func TestUser_GetByEmail(t *testing.T) {
	ur, or, tm := mocks(t)
	userService := NewUser(ur, or, tm)

	user := entity.User{
		ID:        vo.UnsafeID(uuid.New()),
		Name:      mo.Some(vo.UnsafeName("test")),
		Password:  vo.UnsafePasswordHash("hash"),
		Email:     vo.UnsafeEmail("example@example.com"),
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		args    dto.GetByEmailInput
		mock    func()
		wantErr bool
	}{
		{
			name: "success - get user by email",
			args: dto.GetByEmailInput{
				Email: user.Email.Value(),
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "fail - invalid email",
			args: dto.GetByEmailInput{
				Email: "invalid",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "fail - user not found",
			args: dto.GetByEmailInput{
				Email: "notfound@example.com",
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(entity.User{}, entityerr.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail - err from repo",
			args: dto.GetByEmailInput{
				Email: "notfound@example.com",
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(entity.User{}, errSome)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.GetByEmail(context.Background(), tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.args.Email, got.Email.Value())
			}
		})
	}
}

func TestUser_VerifyCred(t *testing.T) {
	ur, or, tm := mocks(t)
	userService := NewUser(ur, or, tm)

	plainPassword, err := vo.NewPlainPassword("password")
	require.NoError(t, err)

	passwordHash, err := vo.NewPasswordHash(plainPassword)
	require.NoError(t, err)

	user := entity.User{
		ID:        vo.UnsafeID(uuid.New()),
		Name:      mo.Some(vo.UnsafeName("test")),
		Password:  passwordHash,
		Email:     vo.UnsafeEmail("example@example.com"),
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		args    dto.VerifyCredInput
		mock    func()
		wantErr bool
	}{
		{
			name: "success - verify cred",
			args: dto.VerifyCredInput{
				Email:    user.Email.Value(),
				Password: plainPassword.Value(),
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "fail - user not found",
			args: dto.VerifyCredInput{
				Email:    "notfound@example.com",
				Password: "password",
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(entity.User{}, entityerr.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail - password not compare",
			args: dto.VerifyCredInput{
				Email:    user.Email.Value(),
				Password: "not compare pass",
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			wantErr: true,
		},
		{
			name: "fail - error from repo",
			args: dto.VerifyCredInput{
				Email:    user.Email.Value(),
				Password: "test",
			},
			mock: func() {
				ur.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(entity.User{}, errSome)
			},
			wantErr: true,
		},
		{
			name: "fail - invalid email",
			args: dto.VerifyCredInput{
				Email:    "invalid",
				Password: "password",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "fail - invalid password",
			args: dto.VerifyCredInput{
				Email:    user.Email.Value(),
				Password: "test",
			},
			mock:    func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := userService.VerifyCred(context.Background(), tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.args.Email, got.Email.Value())
				require.True(t, got.Password.Compare(vo.UnsafePlainPassword(tt.args.Password)))
			}
		})
	}
}
