package paginate

const DefaultLimit = 15

type Paginate interface {
	Limit() int
	Cursor() string
	Token() string
	GenerateToken(string) string
}
