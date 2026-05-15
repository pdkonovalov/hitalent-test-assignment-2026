package postgres

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			return fn(contextWithTx(ctx, tx))
		},
	)
}

func (m *TxManager) DoSerializable(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			return fn(contextWithTx(ctx, tx))
		},
		&sql.TxOptions{Isolation: sql.LevelSerializable},
	)
}

type txContextKey struct{}

func contextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func txFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txContextKey{}).(*gorm.DB)

	return tx, ok
}

type txContextGetter func(context.Context) *gorm.DB

func newTxContextGetter(db *gorm.DB) txContextGetter {
	return func(ctx context.Context) *gorm.DB {
		if tx, ok := txFromContext(ctx); ok {
			return tx
		}

		return db
	}
}
