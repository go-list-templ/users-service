package external

import (
	"encoding/json"
	"fmt"
	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/external/dao"
	"github.com/go-list-templ/grpc/pkg/httpclient"
	"go.uber.org/zap"
)

type UserUnavatar struct {
	client *httpclient.Client
	logger *zap.Logger
}

func NewUserUnavatar(c *httpclient.Client, l *zap.Logger) *UserUnavatar {
	return &UserUnavatar{client: c, logger: l}
}

func (u UserUnavatar) Set(user entity.User) {
	uri := fmt.Sprintf("https://unavatar.io/%v?json", user.Email.Value())

	res, err := u.client.Get(uri)
	if err != nil {
		u.logger.Warn("client get", zap.Error(err))

		return
	}
	defer u.client.ReleaseGet(res)

	avatar := &dao.UserAvatar{}

	err = json.Unmarshal(res.Body(), avatar)
	if err != nil {
		u.logger.Warn("unmarshal", zap.Error(err))

		return
	}

	defer u.client.ReleaseGet(res)

	u.logger.Info("avatar", zap.String("response", avatar.URL))
}
