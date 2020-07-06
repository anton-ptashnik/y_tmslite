package persistence

type TasksRepo struct{}

func (r *TasksRepo) Add(t Task) (int64, error) {
	q := `INSERT INTO tasks (name,description,status_id,project_id,priority_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`

	var id int64
	err := db.QueryRow(q, t.Name, t.Description, t.StatusID, t.ProjectID, t.PriorityID).Scan(&id)
	return id, err
}

func (r *TasksRepo) Get(id int64, pid int64) (Task, error) {
	q := `SELECT * FROM tasks WHERE id=$1`

	var task Task
	err := db.QueryRow(q, id).Scan(&task.ID, &task.StatusID, &task.ProjectID, &task.PriorityID, &task.Name, &task.Description)
	return task, err

}

func (r *TasksRepo) List(pid int64) ([]Task, error) {
	q := `SELECT * FROM tasks WHERE project_id=$1`

	rows, err := db.Query(q, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []Task
	var t Task
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.StatusID, &t.ProjectID, &t.PriorityID, &t.Name, &t.Description)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TasksRepo) Del(id int64, pid int64) error {
	q := `DELETE FROM tasks WHERE id=$1`

	res, err := db.Exec(q, id)
	return verifyModified(res, err)
}

func (r *TasksRepo) Upd(t Task) error {
	q := `UPDATE tasks SET name=$1,description=$2,status_id=$3,priority_id=$4 WHERE id=$5`

	res, err := db.Exec(q, t.Name, t.Description, t.StatusID, t.PriorityID, t.ID)
	return verifyModified(res, err)
}
