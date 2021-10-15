package mysql

import (
	"context"

	"gorm.io/gorm"
)

type CtxTransactionKey struct{}

func CtxWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, CtxTransactionKey{}, tx)
}

type TxImpl struct {
	db *gorm.DB
}

func NewTxImpl(dbcore *Core) *TxImpl {
	return &TxImpl{
		db: dbcore.db,
	}
}

func (t *TxImpl) Transaction(ctx context.Context, fn func(txctx context.Context) error) error {
	db := t.db.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txctx := CtxWithTransaction(ctx, tx)
		return fn(txctx)
	})
}
