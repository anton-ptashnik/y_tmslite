package persistence

func AddTaskStatus(s Status) (int64, error) {
	db := dbConn()
	defer db.Close()
	q := "INSERT INTO statuses (pid, seqNo, name) VALUES ($1,$2,$3) RETURNING id"
	var id int64
	err := db.QueryRow(q, s.PID, s.SeqNo, s.Name).Scan(&id)
	return id, err
}

func DelTaskStatus(id int64) error {
	db := dbConn()
	defer db.Close()
	query := `DELETE FROM statuses WHERE id=$1`
	return verifyModified(db.Exec(query, id))
}

func GetTaskStatus(id int64) (Status, error) {
	db := dbConn()
	defer db.Close()
	q := "SELECT * FROM statuses WHERE id=$1"
	var d Status
	err := db.QueryRow(q, id).Scan(&d.ID, &d.PID, &d.SeqNo, &d.Name)
	return d, err
}

func ListStatuses(projectID int64) ([]Status, error) {
	db := dbConn()
	defer db.Close()
	q := "SELECT * FROM statuses WHERE pid=$1"
	rows, err := db.Query(q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Status
	var d Status
	for rows.Next() {
		rows.Scan(&d.ID, &d.PID, &d.SeqNo, &d.Name)
		res = append(res, d)
	}
	return res, nil
}

func UpdStatus(s Status) error {
	db := dbConn()
	defer db.Close()
	query := `UPDATE statuses SET seqNo=$2,name=$3 WHERE id=$1`
	return verifyModified(db.Exec(query, s.ID, s.SeqNo, s.Name))
}
