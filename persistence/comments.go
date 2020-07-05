package persistence

func AddComment(c Comment) (int64, error) {
	db := dbConn()
	defer db.Close()
	q := `INSERT INTO comments (text, task_id) VALUES ($1,$2) RETURNING id`
	var id int64
	err := db.QueryRow(q, c.Text, c.TaskID).Scan(&id)
	return id, err
}

func GetComment(id int64) (Comment, error) {
	db := dbConn()
	defer db.Close()

	var c Comment
	row := db.QueryRow(`SELECT * FROM comments WHERE id=$1`, id)
	err := row.Scan(&c.ID, &c.TaskID, &c.Text, &c.Date)
	return c, err
}

func ListComments(taskID int64) ([]Comment, error) {
	db := dbConn()
	defer db.Close()
	q := `SELECT * FROM comments WHERE task_id=$1`
	rows, _ := db.Query(q, taskID)
	var comments []Comment
	var c Comment
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.TaskID, &c.Text, &c.Date)
		if err != nil {
			return comments, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func DelComment(id int64) error {
	db := dbConn()
	defer db.Close()
	res, err := db.Exec(`DELETE FROM comments WHERE id=$1`, id)
	return verifyModified(res, err)
}

func UpdComment(c Comment) error {
	db := dbConn()
	defer db.Close()
	q := `UPDATE comments SET text=$1,modified=DEFAULT WHERE id=$2`
	res, err := db.Exec(q, c.Text, c.ID)
	return verifyModified(res, err)
}
