package external

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/external/dao"
	"github.com/go-list-templ/grpc/pkg/httpclient"
)

type UserUnavatar struct {
	client *httpclient.Client
	logger *zap.Logger
}

func NewUserUnavatar(c *httpclient.Client, l *zap.Logger) *UserUnavatar {
	return &UserUnavatar{client: c, logger: l}
}

func (u UserUnavatar) Set(user entity.User) {
	uri := fmt.Sprintf("https://unavatar.io/%v?json", url.QueryEscape(user.Email.Value()))

	res, err := u.client.Get(uri)
	if err != nil {
		u.logger.Info("client uri", zap.String("uri", uri))
		u.logger.Warn("client get", zap.Error(err))
		return
	}
	defer u.client.ReleaseGet(res)

	if res.StatusCode() != http.StatusOK {
		u.logger.Warn("status != 200", zap.Int("status is", res.StatusCode()))
		u.logger.Warn("response", zap.Any("body", res.Body()))
		return
	}

	avatar := &dao.UserAvatar{}

	err = json.Unmarshal(res.Body(), avatar)
	if err != nil {
		u.logger.Warn("unmarshal", zap.Error(err))
		return
	}

	u.logger.Info("user avatar", zap.Any("avatar", avatar))
}
