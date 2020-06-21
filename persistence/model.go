package persistence

type Comment struct {
	taskID uint
	text   string
	date   string
}

type TaskStatus struct {
	projID uint
	name   string
	seqNo  uint
}

type Priority struct {
	name string
}

type Task struct {
	ID                int64
	projID            uint
	name, description string
	status            TaskStatus
	priority          Priority
	comments          []Comment
}

type Project struct {
	ID                int64
	name, description string
	tasks             []Task
	taskStatuses      []TaskStatus
}
