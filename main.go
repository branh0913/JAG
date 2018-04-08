package main

import (
	"JAG/commands"
	"JAG/startup"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"JAG/jobbuilder"
)

type configimpl interface {
	newConfig(configpath string) (*JenkinsConfig, error)
}

type Config struct {
	Endpoint      string            `json:"endpoint"`
	Function      string            `json:"function"`
	LockFile      string            `json:"lockfile"`
	AdminUser     string            `json:"adminuser"`
	AdminPassword string            `json:"adminpassword"`
	PluginScript  map[string]string `json:"pluginscript"`
	ServiceAccount string           `json:"serviceuser"`
	ServicePass    string           `json:"servicepass"`
}

type JenkinsConfig struct {
	Jenkins Config `json:"jenkins"`
}

func (c *JenkinsConfig) newConfig(configpath string) (*JenkinsConfig, error) {

	getconfig := func() *JenkinsConfig {
		jenkinsFile, err := ioutil.ReadFile(configpath)
		if err != nil {
			panic("File cannot be open")
			os.Exit(1)
		}
		json.Unmarshal(jenkinsFile, &c)
		return c
	}()

	return getconfig, nil
}

func check(endpoint, resource, lockfile string, plugins map[string]string) (bool, error) {

	jenkinsInit, err := startup.Init(endpoint, resource, lockfile, plugins)

	if err != nil {
		log.Fatalf("Check failed because %v\n", err)
	}

	if !jenkinsInit {
		fmt.Println("Re-Run JAG once plugins have taken effect on Jenkins")
		return false, nil

	}

	return true, nil

}

func main() {
	config := new(JenkinsConfig)

	configobj, err := config.newConfig("Config/jenkins_automation.json")
	if err != nil {
		log.Fatalf("Cannot load configuation file %v \n", err)
	}

	endpoint := configobj.Jenkins.Endpoint
	function := configobj.Jenkins.Function
	plugins := configobj.Jenkins.PluginScript
	username := configobj.Jenkins.AdminUser
	password := configobj.Jenkins.AdminPassword
	lockfile := configobj.Jenkins.LockFile
	serviceaccount := configobj.Jenkins.ServiceAccount
	servicepass := configobj.Jenkins.ServicePass

	lockfilecheck, err := check(endpoint, function, lockfile, plugins)

	if err != nil {
		log.Fatalf("Something occured when checking file %s\n", err)
	}

	if !lockfilecheck {
		fmt.Println("Check lockfile and rerun jag command")
		os.Exit(1)
	}

	//user create commands
	usercmd := flag.NewFlagSet("user", flag.ExitOnError)
	usercreate := usercmd.Bool("create", false, "Do you want to create user based on Config?")
	unsecure := usercmd.Bool("secure", false, "Create user by authenticating to jenkins?")

	//generate Config commands

	generatecmd := flag.NewFlagSet("generate", flag.ExitOnError)
	generateconfig := generatecmd.String("config", "", "File for generated output")
	generateuser := generatecmd.String("user", "", "User used in Config")

	//CredentialId commands

	credentialcmd := flag.NewFlagSet("credential", flag.ExitOnError)
	credentialCreate := credentialcmd.String("create", "", "Create id of credentialId")
	credentialList := credentialcmd.Bool("list", false, "List credential ids")

	if len(os.Args) < 2 {
		fmt.Println("Type jag user <command>")
		fmt.Println("or jag generate <command>")

	}

	switch os.Args[1] {
	case "user":
		usercmd.Parse(os.Args[2:])
		UserCommands(usercmd, usercreate, unsecure, endpoint, function, username, password)
	case "generate":
		generatecmd.Parse(os.Args[2:])
		GenerateCommands(generatecmd, generateconfig, generateuser, endpoint, function, username, password)
	case "credential":
		credentialcmd.Parse(os.Args[2:])
		CredentialCommand(endpoint, function, username, password, serviceaccount,
			              servicepass, credentialcmd, credentialCreate,
			              credentialList)
	}

}
func CredentialCommand(endpoint, function, username, password, serviceaccount, servicepass string,
	credentialcmd *flag.FlagSet, credentialCreate *string, credentailList *bool) {

		if credentialcmd.Parsed(){
			if *credentialCreate == ""{
				fmt.Println("Please enter a credentialid")
			}
			jbinst := new(jobbuilder.CredentialID)
			credentials, err := jbinst.New(endpoint, function,
				username, password,
				serviceaccount, servicepass,
				credentialCreate)

			if err != nil{
				log.Fatalf("Could not instantiate credential object: %v", err)
			}

			createcredential, err := credentials.Create()

			if err != nil{
				log.Fatalf("Could not create credential id: %v \n", err)
			}

			fmt.Println(createcredential)

			if err != nil{
				log.Fatalf("Could not instantiate new credential object")
			}

		}

}

func UserCommands(usercmd *flag.FlagSet, usercreate *bool,
	              unsecure *bool, endpoint, function,
	              username string, password string) {
	if usercmd.Parsed() {
		if *usercreate && *unsecure {
			fmt.Println("Please enter a valid user command")
			flag.PrintDefaults()

		}
		fmt.Println("Creating new user in jenkins...")
		juser := new(commands.Admin)

		insecureinit, err := juser.New(endpoint, function, username, password)

		if err != nil{
			log.Fatalf("Could not invoke plugin build %v \n", insecureinit)
		}

		insecureinit.Create()


	}
}

func GenerateCommands(generatecmds *flag.FlagSet,
	                  generateconfig, generateuser *string,
	                  endpoint, resource, username, password string){

	   if generatecmds.Parsed(){
	   		if *generateconfig == ""{
	   			fmt.Println("Please supply value for Config")
			}
			if *generateuser == ""{
				fmt.Println("Please suply value for user")
			}
			generate := new(jobbuilder.Token)
			gennew, err := generate.New(endpoint, resource, username, password)

			if err != nil{
				log.Fatalf("Failed to instantiate generate object %v", err)
			}

			apitoken :=  gennew.Retrieve(*generateuser)
			geninst := new(jobbuilder.JBuilderConfig)
			configbuild, err := geninst.New(endpoint, apitoken, *generateconfig)

			if err != nil{
				log.Fatalf("Failed to instantiate new config object %v \n", err.Error())
			}

			buildhandler := configbuild.BuildFile(*generateuser, "templates/jenkins.ini.gotmpl")

			if buildhandler{
				fmt.Printf("jenkins.ini has been generated at %v\n", generateconfig)
			}

	   }

}


