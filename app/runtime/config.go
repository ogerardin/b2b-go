package runtime

type Configuration struct {
	LogLevel    string `short:"l" description:"Log level"`
	Port        int    `short:"p" description:"Port number for web interface"`
	HideConsole bool
	CpuProfile  bool `description:"Activate CPU profiling"`
	MongoPort   int  `short:"mp" description:"Port number for mongod (0 = automatic)"`
}

var CurrentConfig = Configuration{
	LogLevel:    "info",
	Port:        8080,
	CpuProfile:  false,
	HideConsole: false,
	MongoPort:   27017,
}

var DefaultPointersConfig = Configuration{}
