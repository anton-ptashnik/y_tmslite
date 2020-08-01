package persistence

import (
	"database/sql"
	"errors"
	"fmt"
)

func verifyModified(r sql.Result, err error) error {
	if err != nil {
		return err
	}
	rowsAffected, _ := r.RowsAffected()
	if rowsAffected == 0 {
		return errNoMatch
	}
	return nil
}

var (
	errNotImpl = errors.New("not implemented")
	errNoMatch = errors.New("no matching entries")
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTx() (*tx, error) {
	t, err := db.Begin()
	res := tx{t}
	return &res, err
}

type opExecutor interface {
	QueryRow(q string, args ...interface{}) *sql.Row
	Query(q string, args ...interface{}) (*sql.Rows, error)
	Exec(q string, args ...interface{}) (sql.Result, error)
}

type Tx interface {
	opExecutor
	Commit() error
	Rollback() error
}

type tx struct {
	*sql.Tx
}

func (f *tx) Commit() error {
	return f.Commit()
}

func (f *tx) Rollback() error {
	return f.Rollback()
}

func TryCommit(tx Tx) error {
	errC := tx.Commit()
	if errC != nil {
		errR := tx.Rollback()
		if errR != nil {
			return errors.New(fmt.Sprint("Tx discard failed:", errC, errR))
		}
	}
	return nil
}

func initCtx(tx Tx) opExecutor {
	var ctx opExecutor
	if tx == nil {
		ctx = db
	} else {
		ctx = tx
	}
	return ctx
}