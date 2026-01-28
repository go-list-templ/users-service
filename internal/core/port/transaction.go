package port

import "context"

//go:generate mockgen -source=transaction.go -destination=../../../test/mocks/mock_trx.go -package=mocks

type (
	TransactionManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
