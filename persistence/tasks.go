package persistence

func AddTask(task Task) (int64, error) {

	return 0, errNotImpl
}

func GetTask(id int64) (Task, error) {
	return Task{}, errNotImpl

}

func ListTasks(p Project) ([]Task, error) {
	return nil, nil
}

func DelTask(id int64) error {
	return errNotImpl
}

func UpdTask(t Task) error {
	return errNotImpl
}
