package runtime

type Configuration struct {
	LogLevel      string `config:"log"`
	Port          int    `config:"port"`
	HideConsole   bool   `config:"hideconsole"`
	CpuProfile    bool   `config:"prof"`
	MongoPort     int    `config:"mongoport"`
	MongoDataPath string `config:"datapath"`
	Profiles      string `config:"profiles"`
	Version       bool   `config:"version"`
}

func DefaultConfig() *Configuration {
	return &Configuration{
		LogLevel:      "info",
		Port:          8080,
		CpuProfile:    false,
		HideConsole:   false,
		MongoPort:     0,
		MongoDataPath: "",
		Version:       false,
		Profiles:      "",
	}
}

func TestConfig() *Configuration {
	return &Configuration{
		LogLevel:      "debug",
		Port:          8080,
		CpuProfile:    true,
		HideConsole:   false,
		MongoPort:     27017,
		MongoDataPath: "",
		Version:       false,
	}
}

func DefaultPointersConfig() *Configuration {
	return &Configuration{}
}
