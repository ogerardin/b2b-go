package runtime

type Configuration struct {
	LogLevel      string `short:"l" description:"Log level"`
	Port          int    `short:"p" description:"Port number for web interface"`
	HideConsole   bool
	CpuProfile    bool   `description:"Activate CPU profiling"`
	MongoPort     int    `short:"mp" description:"Port number for mongod (0 = automatic)"`
	MongoDataPath string `description:"data path for MongoDB (empty = temp path)"`
	Profiles      string `description:"Active profiles (comma separated)"`
	Version       bool   `description:"Print version information and quits"`
}

var CurrentConfig = defaultConfig()

func defaultConfig() Configuration {
	return Configuration{
		LogLevel:      "info",
		Port:          8080,
		CpuProfile:    false,
		HideConsole:   false,
		MongoPort:     0,
		Version:       false,
		Profiles:      "",
		MongoDataPath: "",
	}
}

var DefaultPointersConfig = Configuration{}
