package main

import (
	"fmt"
	"github.com/brharrelldev/jag/docker"
	"github.com/brharrelldev/jag/jobbuilder"
	"github.com/brharrelldev/jag/startup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var fr FirstRun

// verbs
var (
	create string
)

// jenkins flags
var (
	config *string
)

// Command definitions
var (
	// root
	rootCmd *cobra.Command

	//action command
	createCmd *cobra.Command
)

func init() {

	// Initialize
	viper.AddConfigPath("config")
	viper.SetConfigName("first_run")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	if err := viper.Unmarshal(&fr); err != nil {
		log.Println(err)
	}

}

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

	if check(fr) {
		log.Println("First run detected, initializing startup ")
		dockr, err := docker.NewDockerBuild(&docker.DockerBuildRequest{
			Name:        "Test",
			ArchiveFile: "/tmp/archive.tar",
			DockerFile:  "Dockerfile",
			Tag:         "test-tag",
		})
		if err != nil{
			log.Fatalf("could not instantiate new docker instance %v", err)
		}
		images, err := dockr.BuildDocker()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(images)
	}

	createCmd = &cobra.Command{
		Use: "create [resource]",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) >= 0 {
				fmt.Println("no resource provided")
				cmd.Usage()
				os.Exit(1)
			}

			fmt.Println(config)
			switch args[0] {
			case "jobs":
				if *config == "" {
					fmt.Println("no config file provided")
					cmd.Usage()
					os.Exit(1)
				}

				jobCreateRequest := &jobbuilder.JobCreateRequest{}

				fmt.Println(jobCreateRequest)

			default:
				fmt.Println("unknown resource")
				os.Exit(1)

			}
		},
	}

	config = createCmd.Flags().StringP("config", "c", "config", "config=<path>")

	rootCmd = &cobra.Command{
		Use: "jag [action] [resource]",
	}

	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println("could not execute command", err)
	}

}
