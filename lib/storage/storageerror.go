package storage

type StorageError struct {
	msg string
}

func (se StorageError) Error() string {
	return se.msg
}

func NewStorageError(msg string) error {
	return &StorageError{msg}
}
