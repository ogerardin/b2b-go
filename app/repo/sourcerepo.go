package repo

import (
	"b2b-go/app/domain"
	"b2b-go/lib/mgorepo"
	"github.com/globalsign/mgo"
)

type SourceRepo interface {
	SaveNew(source domain.BackupSource) (string, error)
	Update(id string, source domain.BackupSource) error
	GetById(id string) (domain.BackupSource, error)
	GetAll() ([]domain.BackupSource, error)
	Delete(id string) error
}

type sourceRepoImpl struct {
	mgorepo.Repo
}

var _ SourceRepo = &sourceRepoImpl{}

func NewSourceRepo(s *mgo.Session) SourceRepo {
	repo := &sourceRepoImpl{
		*mgorepo.New(s, "", "sources"),
	}

	//FIXME for testing, remove
	repo.SaveNew(&domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Enabled: true,
			Name:    "my source",
		},
		Paths: []string{"/tmp"},
	})

	return repo
}

func (r *sourceRepoImpl) SaveNew(source domain.BackupSource) (string, error) {
	id, err := r.Repo.SaveNew(&source)
	return id.(string), err
}

func (r *sourceRepoImpl) Update(id string, source domain.BackupSource) error {
	return r.Repo.Update(id, source)
}

func (r *sourceRepoImpl) GetById(id string) (domain.BackupSource, error) {
	retrieved, err := r.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return retrieved.(domain.BackupSource), nil
}

func (r *sourceRepoImpl) GetAll() ([]domain.BackupSource, error) {
	var result []domain.BackupSource
	err := r.Repo.GetAll(&result)
	return result, err
}

func (r *sourceRepoImpl) Delete(id string) error {
	return r.Repo.Delete(id)
}
