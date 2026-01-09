package external

import (
	"encoding/json"
	"fmt"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/external/dao"
	"github.com/go-list-templ/grpc/pkg/http/client"
	"go.uber.org/zap"
)

type UserUnavatar struct {
	client *client.Client
	logger *zap.Logger
}

func NewUserUnavatar(c *client.Client, l *zap.Logger) *UserUnavatar {
	return &UserUnavatar{client: c, logger: l}
}

func (u UserUnavatar) Set(user entity.User) entity.User {
	uri := fmt.Sprintf("https://unavatar.io/%v?json", user.Email.Value())

	res, err := u.client.Get(uri)
	if err != nil {
		u.logger.Warn("client get", zap.Error(err))

		return user
	}
	defer u.client.ReleaseGet(res)

	avatar := &dao.UserAvatar{}

	err = json.Unmarshal(res.Body(), avatar)
	if err != nil {
		u.logger.Warn("unmarshal", zap.Error(err))

		return user
	}

	u.logger.Info("url avatar", zap.String("url", avatar.URL))

	err = user.Avatar.Update(avatar.URL)
	if err != nil {
		u.logger.Warn("avatar set", zap.Error(err))

		return user
	}

	u.logger.Info("updated avatar", zap.Any("avatar", user.Avatar.Value()))

	return user
}
