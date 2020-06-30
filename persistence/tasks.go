package persistence

func AddTask(t Task) (int64, error) {
	q := `INSERT INTO tasks (name,description,status_id,project_id,priority_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	db := dbConn()
	defer db.Close()
	var id int64
	err := db.QueryRow(q, t.Name, t.Description, t.StatusID, t.ProjectID, t.PriorityID).Scan(&id)
	return id, err
}

func GetTask(id int64) (Task, error) {
	q := `SELECT * FROM tasks WHERE id=$1`
	db := dbConn()
	defer db.Close()
	var task Task
	err := db.QueryRow(q, id).Scan(&task.ID, &task.StatusID, &task.ProjectID, &task.PriorityID, &task.Name, &task.Description)
	return task, err

}

func ListTasks(p Project) ([]Task, error) {
	q := `SELECT * FROM tasks WHERE project_id=$1`
	db := dbConn()
	defer db.Close()
	rows, err := db.Query(q, p.ID)
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

func DelTask(id int64) error {
	q := `DELETE FROM tasks WHERE id=$1`
	db := dbConn()
	defer db.Close()
	res, err := db.Exec(q, id)
	return verifyModified(res, err)
}

func UpdTask(t Task) error {
	q := `UPDATE tasks SET name=$1,description=$2,status_id=$3,priority_id=$4 WHERE id=$5`
	db := dbConn()
	defer db.Close()
	res, err := db.Exec(q, t.Name, t.Description, t.StatusID, t.PriorityID, t.ID)
	return verifyModified(res, err)
}
