package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/tkanos/gonfig"
)

// Configuration ...
type Configuration struct {
	ActivePort    string
	PassivePort   string
	JXURL         string
	JXCallbackURL string
}

var configurations Configuration
var once sync.Once

// New ... Set all application configurations here
func New() Configuration {

	once.Do(func() {

		err := gonfig.GetConf(getFileName(), &configurations)
		if err != nil {
			fmt.Println(err)
			os.Exit(500)
		}
		fmt.Println(configurations.ActivePort)

	})
	return configurations
}

func getConfig() Configuration {
	configurations := Configuration{}
	err := gonfig.GetConf("config.developemnt.json", &configurations)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
	fmt.Println(configurations.ActivePort)
	return configurations
}

func getFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	filename := []string{"config.", env, ".json"}
	filePath := path.Join("./config/", strings.Join(filename, ""))

	return filePath
}
