package repo

import (
	"b2b-go/lib/domain"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type BackupSourceRepo struct {
	s    *mgo.Session
	coll string
}

const (
	defaultColl = "sources"
)

func New(s *mgo.Session) *BackupSourceRepo {
	repo := BackupSourceRepo{
		s:    s,
		coll: defaultColl,
	}
	return &repo
}

func (bsr *BackupSourceRepo) Save(s *domain.BackupSource) error {
	session := bsr.s.Copy()
	defer session.Close()

	c := session.DB("").C(bsr.coll)
	_, err := c.UpsertId(s.Id, s)
	return err
}

func (bsr *BackupSourceRepo) GetById(id uint) (*domain.BackupSource, error) {
	// session copy : connection pool
	session := bsr.s.Copy()
	defer session.Close()

	source := domain.BackupSource{}
	c := session.DB("").C(bsr.coll)
	err := c.Find(bson.M{"_id": id}).One(&source)
	return &source, err
}
