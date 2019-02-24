package main

import (
	"fmt"
	"github.com/brharrelldev/jag/docker"
	"github.com/brharrelldev/jag/startup"
	"github.com/spf13/viper"
	"log"
)


type FirstRun struct {
	Dir      string `json:"cachelocation"`
	FileName string `json:"filename"`
}

func check(fr FirstRun) bool {
	path := fr.Dir
	fileName := fr.FileName
	jInit, err := startup.Check(path, fileName)
	if err != nil {
		log.Println("Error occured when trying to check, may not be first run ", err)
		return false
	}

	if !jInit {
		return false
	}

	return true
}


func main() {


	firstRun := FirstRun{}

	viper.AddConfigPath("config")
	viper.SetConfigName("first_run")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	if err := viper.Unmarshal(&firstRun); err != nil {
		log.Println(err)
	}

	fmt.Println(firstRun)

	if check(firstRun) {
		log.Println("First run detected, initializing startup ")
		images, err := docker.BuildDocker()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(images)
	}
}
