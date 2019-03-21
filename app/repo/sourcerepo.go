package repo

import (
	"b2b-go/app/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type sourceRepoImpl struct {
	mgorepo.Repo
}

var _ SourceRepo = &sourceRepoImpl{}

type SourceRepo interface {
	SaveNew(source domain.BackupSource) (bson.ObjectId, error)
	Update(id bson.ObjectId, source domain.BackupSource) error
	GetById(id bson.ObjectId) (domain.BackupSource, error)
	GetAll() ([]domain.BackupSource, error)
	Delete(id bson.ObjectId) error
}

func NewSourceRepo(s *mgo.Session) SourceRepo {
	repo := &sourceRepoImpl{
		*mgorepo.New(s, "sources"),
	}

	//FIXME for testing, remove
	repo.SaveNew(domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Enabled: true,
			Name:    "my source",
		},
		Paths: []string{"/tmp"},
	})

	return repo
}

func (r *sourceRepoImpl) SaveNew(source domain.BackupSource) (bson.ObjectId, error) {
	saved, err := r.Repo.SaveNew(source)
	return saved, err
}

func (r *sourceRepoImpl) Update(id bson.ObjectId, source domain.BackupSource) error {
	return r.Repo.Update(id, source)
}

func (r *sourceRepoImpl) GetById(id bson.ObjectId) (domain.BackupSource, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(domain.BackupSource), err
}

func (r *sourceRepoImpl) GetAll() ([]domain.BackupSource, error) {
	result := make([]domain.BackupSource, 0)
	err := r.Repo.GetAll(&result)
	return result, err
}
