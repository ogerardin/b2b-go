package runtime

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel      string `short:"l" description:"Log level"`
	Port          int    `short:"p" description:"Port number for web interface"`
	HideConsole   bool
	CpuProfile    bool   `description:"Activate CPU profiling"`
	MongoPort     int    `short:"mp" description:"Port number for mongod (0 = automatic)"`
	MongoDataPath string `description:"data path for MongoDB (empty = temporary path)"`
	Profiles      string `description:"Active profiles (comma separated)"`
	Version       bool   `description:"Print version information and quits"`
	Viper         *viper.Viper
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
