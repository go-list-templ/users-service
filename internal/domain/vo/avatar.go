package vo

const FallBackAvatar = "https://unavatar.io/fallback.png"

type Avatar struct {
	value string
}

func NewAvatar() Avatar {
	return Avatar{value: FallBackAvatar}
}

func UnsafeAvatar(avatar string) Avatar {
	return Avatar{value: avatar}
}

func (a Avatar) Value() string {
	return a.value
}
