package persistence

type Comment struct {
	ID     int64
	TaskID int64
	Text   string
	Date   string
}

type Status struct {
	ID    int64
	PID   int64
	Name  string
	SeqNo int
}

type Priority struct {
	ID   int64
	Name string
}

type Task struct {
	ID                int64
	ProjectID         int64
	StatusID          int64
	PriorityID        int64
	Name, Description string
}

type Project struct {
	ID                int64
	Name, Description string
}
