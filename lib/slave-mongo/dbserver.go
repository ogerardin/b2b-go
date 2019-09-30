package slave_mongo

import (
	"bytes"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"gopkg.in/tomb.v2"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// DBServer controls a MongoDB server process to be used within test suites.
//
// The test server is started when Session is called the first time and should
// remain running for the duration of all tests, with the Wipe method being
// called between tests (before each of them) to clear stored data. After all tests
// are done, the Stop method should be called to stop the test server.
type DBServer struct {
	session    *mgo.Session
	output     bytes.Buffer
	server     *exec.Cmd
	dbpath     string
	host       string
	port       int
	logAdapter io.Writer
	tomb       tomb.Tomb
}

func (dbs *DBServer) SetLogAdapter(logAdapter io.Writer) {
	dbs.logAdapter = logAdapter
}

var terminateProcess = func(p *os.Process) {
	p.Signal(os.Interrupt)
}

// SetPath defines the path to the directory where the database files will be
// stored if it is started. The directory path itself is not created or removed
// by the test helper.
func (dbs *DBServer) SetPath(dbpath string) {
	dbs.dbpath = dbpath
}

func (dbs *DBServer) SetPort(port int) {
	dbs.port = port
}

func (dbs *DBServer) Start() {
	if dbs.server != nil {
		//log.Print("DBServer already started")
		return
	}
	//if dbs.dbpath == "" {
	//	panic("DBServer.SetPath must be called before using the server")
	//}
	mgo.SetStats(true)

	// if SetPort has been called, we'll try to listen on the specified port
	// otherwise dbs.port will be 0 and net.Listen will choose a free port that we'll obtain using l.Addr()
	// and pass to mongod
	listenspec := "127.0.0.1:" + strconv.Itoa(dbs.port)
	l, err := net.Listen("tcp", listenspec)
	if err != nil {
		panic("unable to listen on local address" + listenspec + ": " + err.Error())
	}
	addr := l.Addr().(*net.TCPAddr)
	l.Close()
	dbs.host = addr.String()

	args := []string{
		"--bind_ip", "127.0.0.1",
		"--port", strconv.Itoa(addr.Port),
		//"--nssize", "1",
		//"--noprealloc",
		//"--smallfiles",
		//"--nojournal",
	}
	if dbs.dbpath != "" {
		args = append(args, "--dbpath", dbs.dbpath)
	}

	dbs.tomb = tomb.Tomb{}
	dbs.server = exec.Command("mongod", args...)
	//dbs.server.Stdout = &dbs.output
	//dbs.server.Stderr = &dbs.output
	if dbs.logAdapter != nil {
		dbs.server.Stdout = dbs.logAdapter
		dbs.server.Stderr = dbs.logAdapter
	} else {
		dbs.server.Stdout = os.Stdout
		dbs.server.Stderr = os.Stderr
	}
	err = dbs.server.Start()
	if err != nil {
		// print error to facilitate troubleshooting as the panic will be caught in a panic handler
		fmt.Fprintf(os.Stderr, "mongod failed to start: %v\n", err)
		panic(err)
	}
	dbs.tomb.Go(dbs.monitor)
	dbs.Wipe()
}

func (dbs *DBServer) monitor() error {
	dbs.server.Process.Wait()
	if dbs.tomb.Alive() {
		// Present some debugging information.
		fmt.Fprintf(os.Stderr, "---- mongod process died unexpectedly:\n")
		//		fmt.Fprintf(os.Stderr, "%s", dbs.output.Bytes())
		fmt.Fprintf(os.Stderr, "---- mongod processes running right now:\n")
		//FIXME this is Unix-only
		cmd := exec.Command("/bin/sh", "-c", "ps auxw | grep mongod")
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		cmd.Run()
		fmt.Fprintf(os.Stderr, "----------------------------------------\n")

		panic("mongod process died unexpectedly")
	}
	return nil
}

// Stop stops the test server process, if it is running.
//
// It's okay to call Stop multiple times. After the test server is
// stopped it cannot be restarted.
//
// All database sessions must be closed before or while the Stop method
// is running. Otherwise Stop will panic after a timeout informing that
// there is a session leak.
func (dbs *DBServer) Stop() {
	if dbs.session != nil {
		dbs.checkSessions()
		if dbs.session != nil {
			dbs.session.Close()
			dbs.session = nil
		}
	}
	if dbs.server != nil {
		dbs.tomb.Kill(nil)
		terminateProcess(dbs.server.Process)

		select {
		case <-dbs.tomb.Dead():
		case <-time.After(5 * time.Second):
			panic("timeout waiting for mongod process to die")
		}
		dbs.server = nil
	}
}

// Session returns a new session to the server. The returned session
// must be closed after the test is done with it.
//
// The first Session obtained from a DBServer will start it.
func (dbs *DBServer) Session() *mgo.Session {
	if dbs.server == nil {
		dbs.Start()
	}
	if dbs.session == nil {
		mgo.ResetStats()
		var err error
		dbs.session, err = mgo.Dial(dbs.host + "/test")
		if err != nil {
			panic(errors.Wrapf(err, "Failed to connect to MongoDB at %s:%d", dbs.host, dbs.port))
		}
	}
	return dbs.session.Copy()
}

// checkSessions ensures all mgo sessions opened were properly closed.
// For slightly faster tests, it may be disabled setting the
// environment variable CHECK_SESSIONS to 0.
func (dbs *DBServer) checkSessions() {
	if check := os.Getenv("CHECK_SESSIONS"); check == "0" || dbs.server == nil || dbs.session == nil {
		return
	}
	dbs.session.Close()
	dbs.session = nil
	for i := 0; i < 100; i++ {
		stats := mgo.GetStats()
		if stats.SocketsInUse == 0 && stats.SocketsAlive == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("There are mgo sessions still alive.")
}

// Wipe drops all created databases and their data.
//
// The MongoDB server remains running if it was prevoiusly running,
// or stopped if it was previously stopped.
//
// All database sessions must be closed before or while the Wipe method
// is running. Otherwise Wipe will panic after a timeout informing that
// there is a session leak.
func (dbs *DBServer) Wipe() {
	if dbs.server == nil || dbs.session == nil {
		return
	}
	dbs.checkSessions()
	sessionUnset := dbs.session == nil
	session := dbs.Session()
	defer session.Close()
	if sessionUnset {
		dbs.session.Close()
		dbs.session = nil
	}
	names, err := session.DatabaseNames()
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		switch name {
		case "admin", "local", "config":
		default:
			err = session.DB(name).DropDatabase()
			if err != nil {
				panic(err)
			}
		}
	}
}
