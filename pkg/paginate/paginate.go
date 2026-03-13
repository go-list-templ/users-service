package paginate

type Paginate interface {
	Limit() int
	Cursor() string
	Token() string
	GenerateToken(string) string
}
