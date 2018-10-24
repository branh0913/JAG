package jobbuilder

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type Config interface {
	New(endpoint, apitoken, jconfigpath string) (*JBuilderConfig, error)
	BuildFile(currentuser string, templatePath string) bool
}

type JBuilderConfig struct {
	APIToken    string
	Endpoint    string
	JConfigPath string
}

func (j *JBuilderConfig) New(endpoint, apitoken, jconfigpath string) (*JBuilderConfig, error) {

	return &JBuilderConfig{
		Endpoint:    endpoint,
		APIToken:    apitoken,
		JConfigPath: jconfigpath,
	}, nil
}

func (j JBuilderConfig) BuildFile(currentuser string, templatePath string) bool {

	type CurrentUser struct {
		CUser    string
		APIToken string
		Endpoint string
	}

	cuserinst := CurrentUser{CUser: currentuser,
		APIToken: strings.Trim(j.APIToken, "\n"),
		Endpoint: strings.Replace(j.Endpoint, "/scriptText", "", 1)}

	t, err := template.ParseFiles(templatePath)

	if err != nil {
		log.Fatalf("Could not generate template because %v\n", err)
		return false
	}

	f, err := os.Create(j.JConfigPath)

	if err != nil {
		log.Fatalf("Could not create file because of %v \n", err)
		os.Exit(2)
		return false
	}

	t.Execute(f, &cuserinst)
	return true
}
