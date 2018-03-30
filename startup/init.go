package startup

import (

	"os"
	"log"
	"JAG/commands"

	"fmt"
)

func Init( endpoint, resource, jaglock string, plugins map[string]string) (bool, error){

	if _, err := os.Stat(jaglock); os.IsNotExist(err){
		log.Printf("File jag.lock is not detected, invoking first run %s", err)
		log.Println("Intializing loading plugins")
		resp := new(commands.PluginInsecure)

		initPlugins, err :=  resp.New(endpoint, resource, plugins)

		if err != nil{
			log.Fatalf("Initialization of Plugin struct failed %v \n", err)
		}

		pluginload, err := initPlugins.Create()

		fmt.Println(pluginload)

		if err != nil{
			log.Fatalf("Loading plugins threw an error %v \n", err)
		}

		log.Println("Creating lockfile...")
		lockfile, err := os.Create(jaglock)

		if err != nil{
			log.Fatalf("Major error, file cannot be created because %v", err)
		}

		log.Printf("File %s has successsfuy been created\n", lockfile.Name())
		log.Printf("Please restart jenkins so that changes take effect\n")
		return false, nil

	}

	log.Println("Ready to start jag")
	return true, nil
}

