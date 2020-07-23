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

func NewTx() (*Tx, error) {
	t, err := db.Begin()
	res := Tx{t}
	return &res, err
}

type opExecutor interface {
	QueryRow(q string, args ...interface{}) *sql.Row
	Query(q string, args ...interface{}) (*sql.Rows, error)
	Exec(q string, args ...interface{}) (sql.Result, error)
}

type Tx struct {
	tx *sql.Tx
}

func (f *Tx) Commit() error {
	return f.tx.Commit()
}

func (f *Tx) Rollback() error {
	return f.tx.Rollback()
}

func TryCommit(tx *Tx) error {
	errC := tx.Commit()
	if errC != nil {
		errR := tx.Rollback()
		if errR != nil {
			return errors.New(fmt.Sprint("Tx discard failed:", errC, errR))
		}
	}
	return nil
}