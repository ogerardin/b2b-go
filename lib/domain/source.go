package domain

type BackupSource struct {
	Id      uint `bson:"_id"`
	Enabled bool
	Name    string
}
