package main

import (
	"fmt"
	"github.com/brharrelldev/jag/startup"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"log"
)


//
//// verbs
//var (
//	create string
//)
//
//// jenkins flags
//var (
//	config *string
//)
//
//// Command definitions
//var (
//	// root
//	rootCmd *cobra.Command
//
//	//action command
//	createCmd *cobra.Command
//)
//

//}

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

func before(appctx *cli.Context) error  {
	var fr FirstRun
	viper.AddConfigPath("config")
	viper.SetConfigName("first_run")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read config for unknow reason %v", err)
	}

	if err := viper.Unmarshal(&fr); err != nil {
		return fmt.Errorf("error unmarshalling json %v", err)
	}

	return nil

}

func main() {

	app := cli.NewApp()
	app.Name = "jag"
	app.Before = before


}
