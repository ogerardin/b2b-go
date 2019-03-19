package repo

import (
	"b2b-go/app"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type sourceRepoImpl struct {
	mgorepo.Repo
}

var _ SourceRepo = &sourceRepoImpl{}

type SourceRepo interface {
	SaveNew(source app.BackupSource) (bson.ObjectId, error)
	Update(id bson.ObjectId, source app.BackupSource) error
	GetById(id bson.ObjectId) (app.BackupSource, error)
	GetAll() ([]app.BackupSource, error)
	Delete(id bson.ObjectId) error
}

func NewSourceRepo(s *mgo.Session) SourceRepo {
	return &sourceRepoImpl{
		*mgorepo.New(s, "sources"),
	}
}

func (r *sourceRepoImpl) SaveNew(source app.BackupSource) (bson.ObjectId, error) {
	saved, err := r.Repo.SaveNew(source)
	return saved, err
}

func (r *sourceRepoImpl) Update(id bson.ObjectId, source app.BackupSource) error {
	return r.Repo.Update(id, source)
}

func (r *sourceRepoImpl) GetById(id bson.ObjectId) (app.BackupSource, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(app.BackupSource), err
}

func (r *sourceRepoImpl) GetAll() ([]app.BackupSource, error) {
	var result []app.BackupSource
	err := r.Repo.GetAll(&result)
	return result, err
}
