package persistence

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

type taskTests struct {
	TasksRepo
}

type taskTestsInput struct {
	p     Project
	s     Status
	pri   Priority
	tasks []Task
}

func TestTasks(t *testing.T) {
	_, err := InitDb()
	panicOnErr(err)
	defer db.Close()

	input, err := prepareTaskTests(db)
	if err != nil {
		t.Fatal(err)
	}

	tests := taskTests{}
	t.Run("get", tests.getTask(input.tasks))
	t.Run("list", tests.listTasks(input.p.ID))
	t.Run("add", tests.addTask(input.p, input.s, input.pri))
	t.Run("upd", tests.updTask(input.tasks[0]))
	t.Run("del", tests.delTask(input.tasks[0]))
}
func (test *taskTests) addTask(p Project, s Status, pri Priority) func(t *testing.T) {
	newTask := Task{
		ProjectID:   p.ID,
		StatusID:    s.ID,
		PriorityID:  pri.ID,
		Name:        "testtask",
		Description: "test purpose entry",
	}
	return func(t *testing.T) {
		taskID, err := test.TasksRepo.Add(newTask)
		if err != nil {
			t.Fatal(err)
		}
		newTask.ID = taskID
		if !checkTaskExists(newTask) {
			t.Error("responded ok but entry not found")
		}
	}
}
func (test *taskTests) getTask(expectedTasks []Task) func(*testing.T) {
	return func(t *testing.T) {
		for _, expectedTask := range expectedTasks {
			actual, err := test.TasksRepo.Get(expectedTask.ID, expectedTask.ProjectID)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expectedTask, actual) {
				t.Error("expected/actual entry mismatch:", expectedTask, actual)
			}
		}
	}
}
func (test *taskTests) listTasks(pid int64) func(t *testing.T) {
	return func(t *testing.T) {
		expected, err := test.TasksRepo.List(pid)
		if err != nil {
			t.Fatal(err)
		}
		for _, task := range expected {
			if !checkTaskExists(task) {
				t.Error("task not found:", task)
			}
		}
	}
}
func (test *taskTests) delTask(task Task) func(t *testing.T) {
	return func(t *testing.T) {
		err := test.TasksRepo.Del(task.ID, task.ProjectID)
		if err != nil {
			t.Fatal(err)
		}
		if checkTaskExists(task) {
			t.Error("responded ok but entry left in DB:", task)
		}
	}

}
func (test *taskTests) updTask(task Task) func(t *testing.T) {
	return func(t *testing.T) {
		updTask := task
		updTask.Name += "_updated"
		err := test.TasksRepo.Upd(updTask)
		if err != nil {
			t.Fatal(err)
		}
		if !checkTaskExists(updTask) {
			t.Error("responded ok but updated entry not found. Old/new:", task, updTask)
		}
	}
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

func checkTaskExists(expected Task) bool {
	q := `SELECT * FROM tasks WHERE id=$1`
	var actual Task
	err := db.QueryRow(q, expected.ID).Scan(&actual.ID, &actual.StatusID, &actual.ProjectID, &actual.PriorityID, &actual.Name, &actual.Description)
	return err == nil && reflect.DeepEqual(expected, actual)
}

func prepareTaskTests(db *sql.DB) (taskTestsInput, error) {
	projects, err := prepareProjects(db, 1)
	if err != nil {
		return taskTestsInput{}, err
	}
	statuses, err := prepareStatuses(db, projects[0], 1)
	if err != nil {
		return taskTestsInput{}, err
	}
	priorities, err := preparePriorities(db, 1)
	if err != nil {
		return taskTestsInput{}, err
	}
	tasks, err := prepareTasks(db, projects[0], statuses[0], priorities[0], 2)
	if err != nil {
		return taskTestsInput{}, err
	}
	return taskTestsInput{
		p:     projects[0],
		s:     statuses[0],
		pri:   priorities[0],
		tasks: tasks,
	}, nil
}
