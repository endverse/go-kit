package mysql

import (
	"context"

	"gorm.io/gorm"
)

type ctxTransactionKey struct{}

func CtxWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxTransactionKey{}, tx)
}

type txImpl struct {
	db *gorm.DB
}

func NewTxImpl(dbcore *core) *txImpl {
	return &txImpl{
		db: dbcore.db,
	}
}

func (t *txImpl) Transaction(ctx context.Context, fn func(txctx context.Context) error) error {
	db := t.db.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txctx := CtxWithTransaction(ctx, tx)
		return fn(txctx)
	})
}
