package repo

import (
	"b2b-go/app/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TargetRepo interface {
	SaveNew(source domain.BackupTarget) (bson.ObjectId, error)
	Update(id bson.ObjectId, target domain.BackupTarget) error
	GetById(id bson.ObjectId) (domain.BackupTarget, error)
	GetAll() ([]domain.BackupTarget, error)
	Delete(id bson.ObjectId) error
}

type targetRepoImpl struct {
	mgorepo.Repo
}

var _ TargetRepo = &targetRepoImpl{}

func (r *targetRepoImpl) Update(id bson.ObjectId, target domain.BackupTarget) error {
	return r.Repo.Update(id, target)
}

func (r *targetRepoImpl) GetAll() ([]domain.BackupTarget, error) {
	var result []domain.BackupTarget
	err := r.Repo.GetAll(result)
	return result, err
}

func NewTargetRepo(s *mgo.Session) TargetRepo {
	return &targetRepoImpl{
		*mgorepo.New(s, "targets"),
	}
}

func (r *targetRepoImpl) SaveNew(target domain.BackupTarget) (bson.ObjectId, error) {
	saved, err := r.Repo.SaveNew(target)
	return saved, err
}

func (r *targetRepoImpl) GetById(id bson.ObjectId) (domain.BackupTarget, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(domain.BackupTarget), err
}
