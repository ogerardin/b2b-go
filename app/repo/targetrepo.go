package repo

import (
	"b2b-go/app/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
)

type TargetRepo interface {
	SaveNew(source domain.BackupTarget) (string, error)
	Update(id string, target domain.BackupTarget) error
	GetById(id string) (domain.BackupTarget, error)
	GetAll() ([]domain.BackupTarget, error)
	Delete(id string) error
}

type targetRepoImpl struct {
	mgorepo.Repo
}

var _ TargetRepo = &targetRepoImpl{}

func (r *targetRepoImpl) Update(id string, target domain.BackupTarget) error {
	return r.Repo.Update(id, target)
}

func (r *targetRepoImpl) GetAll() ([]domain.BackupTarget, error) {
	var result []domain.BackupTarget
	err := r.Repo.GetAll(result)
	return result, err
}

func NewTargetRepo(s *mgo.Session) TargetRepo {
	return &targetRepoImpl{
		*mgorepo.New(s, "", "targets"),
	}
}

func (r *targetRepoImpl) SaveNew(target domain.BackupTarget) (string, error) {
	id, err := r.Repo.SaveNew(target)
	return id.(string), err
}

func (r *targetRepoImpl) GetById(id string) (domain.BackupTarget, error) {
	retrieved, err := r.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return retrieved.(domain.BackupTarget), nil
}

func (r *targetRepoImpl) Delete(id string) error {
	return r.Repo.Delete(id)
}
