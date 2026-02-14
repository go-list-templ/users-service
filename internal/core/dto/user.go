package dto

type (
	UserCreateInput struct {
		Name  string
		Email string
	}

	UserListInput struct {
		PageSize  int64
		PageToken string
	}
)
