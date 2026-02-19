package paginate

const (
	DefaultLimitList = 10
	MaxLimitList     = 100
)

type Paginate interface {
	Limit() int
	Cursor() string
	GenerateToken(string) string
}
