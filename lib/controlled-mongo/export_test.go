package controlled_mongo

import (
	"os"
)

func (dbs *DBServer) ProcessTest() *os.Process {
	if dbs.server == nil {
		return nil
	}
	return dbs.server.Process
}
