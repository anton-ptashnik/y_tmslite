package persistence

type Comment struct {
	ID     int64
	taskID uint
	text   string
	date   string
}

type TaskStatus struct {
	ID     int64
	projID uint
	name   string
	seqNo  uint
}

type Priority struct {
	ID   int64
	Name string
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
	Name, Description string
	Tasks             []Task       `json:",omitempty"`
	TaskStatuses      []TaskStatus `json:",omitempty"`
}
