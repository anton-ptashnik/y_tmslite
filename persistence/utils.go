package persistence

import "database/sql"

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
