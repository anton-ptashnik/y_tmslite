package persistence

func AddProject(p Project) (int64, error) {
	q := `INSERT INTO projects (name, description) VALUES ($1,$2) RETURNING id`
	var id int64
	err := db.QueryRow(q, p.Name, p.Description).Scan(&id)
	return id, err
}

func ListProjects() ([]Project, error) {
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
	var p Project
	err := db.QueryRow(`SELECT * FROM projects WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Description)
	return p, err
}

func UpdProject(p Project) error {
	query := `UPDATE projects SET name=$2, description=$3 WHERE id=$1`
	return verifyModified(db.Exec(query, p.ID, p.Name, p.Description))
}

func DelProject(id int64) error {
	query := `DELETE FROM projects WHERE id=$1`
	return verifyModified(db.Exec(query, id))
}
