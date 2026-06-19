package rcpostgres

import (
	"context"

	"gorm.io/gorm"
)

type contextKeyTx struct{}

func ctxWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKeyTx{}, tx)
}

func getTxFromCtx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(contextKeyTx{}).(*gorm.DB)
	return tx, ok
}
