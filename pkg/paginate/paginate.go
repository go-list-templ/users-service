package paginate

const DefaultLimit = 15

type Paginate interface {
	Limit() int
	Cursor() string
	GenerateToken(string) string
}
