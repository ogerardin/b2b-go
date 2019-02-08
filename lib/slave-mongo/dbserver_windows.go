package slave_mongo

import (
	"github.com/mattn/psutil"
	"os"
	"syscall"
)

func init() {
	terminateProcess = func(p *os.Process) {
		psutil.TerminateTree(p.Pid, int(syscall.SIGKILL))
	}
}
