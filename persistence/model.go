package persistence

type Comment struct {
	taskID uint
	text   string
	date   string
}

type Status struct {
	name string
	pos  uint
}

type Priority struct {
	name string
}

type Task struct {
	projID            uint
	name, description string
	status            Status
	priority          Priority
	comments          []Comment
}

type Project struct {
	name, description string
	tasks             []Task
}
