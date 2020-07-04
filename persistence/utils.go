package persistence

import (
	"database/sql"
	"errors"
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
