package service

import "y_finalproject/persistence"

type ProjectsRepo interface {
	Add(project persistence.Project) (int64,error)
	Get(id int64) (persistence.Project,error)
	Del(id int64) error
	Upd(project persistence.Project) error
	List() ([]persistence.Project, error)
}

type InsertStatusOp func(status persistence.Status) (int64, error)
type InsertStatusOpTx func(tx persistence.Tx) InsertStatusOp

type ProjectsRepoTx func(tx persistence.Tx) ProjectsRepo

type ProjectsService struct {
	ProjectsRepo
	ProjectsRepoTx
	InsertStatusOpTx
	TxInitiator
}

func (s *ProjectsService) Add(p persistence.Project) (int64,error)  {
	tx, _ := s.TxInitiator()
	r := s.ProjectsRepoTx(tx)
	pid, _ := r.Add(p)
	s.initStatus(tx, pid)
	return pid, persistence.TryCommit(tx)
}

func (s *ProjectsService) initStatus(tx persistence.Tx, pid int64) error {
	op := s.InsertStatusOpTx(tx)
	_, err := op(persistence.Status{
		PID:   pid,
		Name:  "default",
		SeqNo: 1,
	})
	return err
}