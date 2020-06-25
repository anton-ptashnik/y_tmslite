package persistence

import (
	"errors"
)

var (
	errNoMatch = errors.New("no matching entries")
)

func AddProject(p Project) (int64, error) {
	db := dbConn()
	defer db.Close()
	q := `INSERT INTO projects (name, description) VALUES ($1,$2) RETURNING id`
	var id int64
	err := db.QueryRow(q, p.Name, p.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ListProjects() ([]Project, error) {
	db := dbConn()
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var p Project
	var res []Project
	for rows.Next() {
		rows.Scan(&p.ID, &p.Name, &p.Description)
		res = append(res, p)
	}
	return res, nil
}
func GetProject(id int64) (Project, error) {
	db := dbConn()
	defer db.Close()
	var p Project
	err := db.QueryRow(`SELECT * FROM projects WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Description)
	return p, err
}

func UpdProject(p Project) error {
	db := dbConn()
	defer db.Close()
	query := `UPDATE projects SET name=$2, description=$3 WHERE id=$1 returning id`
	res, err := db.Exec(query, p.ID, p.Name, p.Description)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errNoMatch
	}
	return nil
}

func DelProject(id int64) error {
	db := dbConn()
	defer db.Close()
	query := `DELETE FROM projects WHERE id=$1 returning id`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return errors.New("no match")
	}
	return nil
}
