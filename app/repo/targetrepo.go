package repo

import (
	"b2b-go/app"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TargetRepo interface {
	SaveNew(source app.BackupTarget) (bson.ObjectId, error)
	GetById(id bson.ObjectId) (app.BackupTarget, error)
}

type targetRepoImpl struct {
	mgorepo.Repo
}

func NewTargetRepo(s *mgo.Session) TargetRepo {
	return &targetRepoImpl{
		*mgorepo.New(s, "targets"),
	}
}

func (r *targetRepoImpl) SaveNew(target app.BackupTarget) (bson.ObjectId, error) {
	saved, err := r.Repo.SaveNew(target)
	return saved, err
}

func (r *targetRepoImpl) GetById(id bson.ObjectId) (app.BackupTarget, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(app.BackupTarget), err
}
