package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
)

type sourceRepoImpl struct {
	mgorepo.Repo
}

type SourceRepo interface {
	SaveNew(source domain.BackupSource) (interface{}, error)
	GetById(id interface{}) (domain.BackupSource, error)
}

func NewSourceRepo(s *mgo.Session) SourceRepo {
	return &sourceRepoImpl{
		*mgorepo.New(s, "sources"),
	}
}

func (r *sourceRepoImpl) SaveNew(source domain.BackupSource) (interface{}, error) {
	saved, err := r.Repo.SaveNew(source)
	return saved, err
}

func (r *sourceRepoImpl) GetById(id interface{}) (domain.BackupSource, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(domain.BackupSource), err
}
