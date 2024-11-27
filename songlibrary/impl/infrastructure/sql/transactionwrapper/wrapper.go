package transactionwrapper

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionWrapper struct {
	conn *pgxpool.Pool
}

func NewTransactionWrapper(conn *pgxpool.Pool) *TransactionWrapper {
	return &TransactionWrapper{conn}
}

func (t *TransactionWrapper) ExecuteWithTransaction(context context.Context, fn func(context.Context, pgx.Tx) error) error {
	tx, err := t.conn.Begin(context)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(context)
		} else {
			err = tx.Commit(context)
		}
	}()

	err = fn(context, tx)
	if err != nil {
		rollbackErr := tx.Rollback(context)
		if rollbackErr != nil {
			return errors.New("rollback error: " + rollbackErr.Error())
		}
		return err
	}

	return nil
}
