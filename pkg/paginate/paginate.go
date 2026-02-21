package paginate

const DefaultLimit = 30

type Paginate interface {
	Limit() int
	Cursor() string
	GenerateToken(string) string
}
