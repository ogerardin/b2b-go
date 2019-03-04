package mongorepository

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/genericrepo"
	"github.com/globalsign/mgo"
)

type SourceRepoImpl struct {
	genericrepo.Repo
}

type SourceRepo interface {
	SaveNew(source domain.BackupSource) (interface{}, error)
	GetById(id interface{}) (domain.BackupSource, error)
}

func NewSourceRepo(s *mgo.Session) *SourceRepoImpl {
	return &SourceRepoImpl{
		*genericrepo.NewRepo(s, "sources"),
	}
}

func (r *SourceRepoImpl) SaveNew(source domain.BackupSource) (interface{}, error) {
	saved, err := r.Repo.SaveNew(source)
	return saved, err
}

func (r *SourceRepoImpl) GetById(id interface{}) (domain.BackupSource, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(domain.BackupSource), err
}
