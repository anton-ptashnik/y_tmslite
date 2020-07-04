package persistence

func AddComment(c Comment) (int64, error) {
	return 0, errNotImpl
}

func GetComment(id int64) (Comment, error) {
	return Comment{}, errNotImpl
}

func ListComments(taskID int64) ([]Comment, error) {
	return nil, errNotImpl
}

func DelComment(id int64) error {
	return nil
}

func UpdComment(c Comment) error {
	return errNotImpl
}
