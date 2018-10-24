package main

import (
	"github.com/brharrelldev/jag/startup"
	"github.com/spf13/viper"
	"log"
)

type FirstRun struct {
	Dir      string `json:"cachelocation"`
	FileName string `json:"filename"`
}
type Config struct {
	Endpoint       string            `json:"endpoint"`
	Function       string            `json:"function"`
	LockFile       string            `json:"lockfile"`
	AdminUser      string            `json:"adminuser"`
	AdminPassword  string            `json:"adminpassword"`
	PluginScript   map[string]string `json:"pluginscript"`
	ServiceAccount string            `json:"serviceuser"`
	ServicePass    string            `json:"servicepass"`
}

//type JenkinsConfig struct {
//	Jenkins Config `json:"jenkins"`
//}

//Wrapper to check if there has been a first run
func check(fr FirstRun) bool {
	path := fr.Dir
	fileName := fr.FileName
	jInit, err := startup.Check(path, fileName)
	if err != nil {
		log.Println("Error occured when trying to check, may not be first run ", err)
		return false
	}

	if jInit {
		return false
	}

	return true
}

func main() {
	firstRun := FirstRun{}

	viper.AddConfigPath("config")
	viper.SetConfigName("first_run")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("could not read config due to ", err)
	}

	viper.Unmarshal(&firstRun)

	if check(firstRun) {
		log.Println("First run detected, initializing startup ")
	}

}
