package persistence

type Comment struct {
	ID     int64
	taskID uint
	text   string
	date   string
}

type TaskStatus struct {
	ID    int64
	PID   int64
	Name  string
	SeqNo int64
}

type Priority struct {
	ID   int64
	Name string
}

type Task struct {
	ID                int64
	PID               uint
	Name, Description string
	Status            TaskStatus
	Priority          Priority
	Comments          []Comment
}

type Project struct {
	ID                int64
	Name, Description string
	Tasks             []Task       `json:",omitempty"`
	TaskStatuses      []TaskStatus `json:",omitempty"`
}
