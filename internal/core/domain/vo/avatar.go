package vo

import (
	"errors"
	"net/url"
)

const FallBackAvatar = "https://unavatar.io/fallback.png"

var ErrInvalidAvatarURL = errors.New("invalid avatar URL")

type Avatar struct {
	value string
}

func NewAvatar() Avatar {
	return Avatar{value: FallBackAvatar}
}

func UnsafeAvatar(avatar string) Avatar {
	return Avatar{value: avatar}
}

func (a *Avatar) Value() string {
	return a.value
}

func (a *Avatar) Update(avatar string) error {
	u, err := url.ParseRequestURI(avatar)
	if err != nil {
		return err
	}

	if u.Scheme == "" && u.Host == "" {
		return ErrInvalidAvatarURL
	}

	a.value = avatar

	return nil
}
