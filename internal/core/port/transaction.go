package port

import "context"

//go:generate mockgen -source=transaction.go -destination=../../../test/mocks/mocks_trx_test.go -package=mock_test

type (
	TransactionManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
