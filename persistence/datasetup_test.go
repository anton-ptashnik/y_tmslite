package persistence

import (
	"database/sql"
	"fmt"
)

func prepareProjects(db *sql.DB, n int) ([]Project, error) {
	var p = Project{
		Name:        "newp",
		Description: "somethingverImportant",
	}
	q := `INSERT INTO projects (name, description) VALUES ($1,$2) RETURNING id`
	var res []Project
	for i := n; i > 0; i-- {
		pc := p
		pc.Name += fmt.Sprint(i)
		err := db.QueryRow(q, pc.Name, pc.Description).Scan(&pc.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, pc)
	}
	return res, nil
}

func prepareStatuses(db *sql.DB, p Project, n int) ([]Status, error) {
	var baseStatus = Status{
		PID:   p.ID,
		Name:  "default",
		SeqNo: 0,
	}
	q := `INSERT INTO statuses (pid, name, seqNo) VALUES ($1,$2,$3) RETURNING id`
	var statuses []Status
	for n > 0 {
		n--
		status := baseStatus
		status.Name = fmt.Sprint(status.Name, n)
		status.SeqNo = n
		err := db.QueryRow(q, status.PID, status.Name, status.SeqNo).Scan(&status.ID)
		if err != nil {
			return statuses, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func preparePriorities(db *sql.DB, n int) ([]Priority, error) {
	q := `INSERT INTO priorities (name) VALUES ($1) RETURNING id`
	var res []Priority
	baseName := "No"
	for n > 0 {
		n--
		p := Priority{
			Name: fmt.Sprint(baseName, n),
		}
		err := db.QueryRow(q, p.Name).Scan(&p.ID)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}
	return res, nil
}

func prepareTasks(db *sql.DB, p Project, s Status, pri Priority, n int) ([]Task, error) {
	baseTask := Task{
		ProjectID:   p.ID,
		StatusID:    s.ID,
		PriorityID:  pri.ID,
		Name:        "tesTask",
		Description: "test purpose entry",
	}
	var tasks []Task
	q := `INSERT INTO tasks (project_id,status_id,priority_id,name,description) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	for n > 0 {
		n--
		t := baseTask
		t.Name = fmt.Sprint(t.Name, n)
		db.QueryRow(q, t.ProjectID, t.StatusID, t.PriorityID, t.Name, t.Description).Scan(&t.ID)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func prepareComments(db *sql.DB, taskID int64, n int) []Comment {
	q := `INSERT INTO comments (text,task_id) VALUES ($1,$2) RETURNING id`
	baseText := "testcomment_"
	var res []Comment
	for n > 0 {
		n--
		comment := Comment{
			TaskID: taskID,
			Text:   fmt.Sprint(baseText, n),
		}
		err := db.QueryRow(q, comment.Text, comment.TaskID).Scan(&comment.ID)
		if err != nil {
			panic(err)
		}
		res = append(res, comment)
	}
	return res
}

