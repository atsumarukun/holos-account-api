package database

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain"
)

type transactionKey struct{}

type transactionObject struct {
	db *sqlx.DB
}

func NewDBTransactionObject(db *sqlx.DB) domain.TransactionObject {
	return &transactionObject{
		db: db,
	}
}

func (to *transactionObject) Transaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := to.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(rollbackErr.Error())
			}
		}
	}()

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	if err := fn(ctx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err.Error())
	}

	return nil
}

type driver interface {
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.Execer
	sqlx.ExecerContext
}

func getDriver(ctx context.Context, db *sqlx.DB) driver {
	if tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
