package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
)

type TargetRepo interface {
	SaveNew(source domain.BackupTarget) (interface{}, error)
	GetById(id interface{}) (domain.BackupTarget, error)
}

type targetRepoImpl struct {
	mgorepo.Repo
}

func NewTargetRepo(s *mgo.Session) TargetRepo {
	return &targetRepoImpl{
		*mgorepo.New(s, "targets"),
	}
}

func (r *targetRepoImpl) SaveNew(target domain.BackupTarget) (interface{}, error) {
	saved, err := r.Repo.SaveNew(target)
	return saved, err
}

func (r *targetRepoImpl) GetById(id interface{}) (domain.BackupTarget, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(domain.BackupTarget), err
}
