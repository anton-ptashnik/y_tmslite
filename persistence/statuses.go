package persistence

type StatusesRepo struct {
	opExecutor
}

func (r *StatusesRepo) Add(s Status) (int64, error) {
	q := "INSERT INTO statuses (pid, seqNo, name) VALUES ($1,$2,$3) RETURNING id"
	var id int64
	err := r.QueryRow(q, s.PID, s.SeqNo, s.Name).Scan(&id)
	return id, err
}
func (r *StatusesRepo) Del(id int64, pid int64) error {
	query := `DELETE FROM statuses WHERE id=$1`
	return verifyModified(r.Exec(query, id))
}

func (r *StatusesRepo) Get(id int64, pid int64) (Status, error) {

	q := "SELECT * FROM statuses WHERE id=$1"
	var d Status
	err := r.QueryRow(q, id).Scan(&d.ID, &d.PID, &d.SeqNo, &d.Name)
	return d, err
}

func (r *StatusesRepo) List(projectID int64) ([]Status, error) {

	q := "SELECT * FROM statuses WHERE pid=$1"
	rows, err := r.Query(q, projectID)
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

func (r *StatusesRepo) Upd(s Status) error {

	query := `UPDATE statuses SET seqNo=$2,name=$3 WHERE id=$1`
	return verifyModified(r.Exec(query, s.ID, s.SeqNo, s.Name))
}

func NewStatusesRepo(tx *Tx) *StatusesRepo {
	var ctx opExecutor
	if tx == nil {
		ctx = db
	} else {
		ctx = tx.tx
	}
	return &StatusesRepo{ctx}
}