package port

import "context"

//go:generate mockgen -source=transaction.go -destination=mock/mock_trx.go -package=mock

type (
	TransactionManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
