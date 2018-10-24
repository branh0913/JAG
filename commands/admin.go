package commands

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type JenkinsUser interface {
	New(endpoint, resource, username, password string) (error, *Admin)
	Create() (string, error)
}

type Admin struct {
	Endpoint string
	Resource string
	Username string
	Password string
}

func (u *Admin) New(endpoint, resource, username, password string) (*Admin, error) {

	return &Admin{
		Endpoint: endpoint,
		Resource: resource,
		Username: username,
		Password: password,
	}, nil
}

func (u Admin) Create() (string, error) {

	userCreateGroovy := `import jenkins.model.*
import hudson.security.*
def env = System.getenv()
def jenkins = Jenkins.getInstance()
jenkins.setSecurityRealm(new HudsonPrivateSecurityRealm(false))
jenkins.setAuthorizationStrategy(new GlobalMatrixAuthorizationStrategy())
def user = jenkins.getSecurityRealm().createAccount("{{.Username}}", "{{.Password}}")
user.save()
jenkins.getAuthorizationStrategy().add(Jenkins.ADMINISTER, "admin")
jenkins.save()
`
	tmpl := template.Must(template.New("CreateUser").Parse(userCreateGroovy))

	var bytetmpl bytes.Buffer
	err := tmpl.Execute(&bytetmpl, u)

	if err != nil {
		log.Fatalf("Issue with template %v", tmpl)
	}

	scriptString := bytetmpl.String()

	scriptExec := func() string {

		data := url.Values{}
		data.Set(u.Resource, string(scriptString))
		req, err := http.PostForm(u.Endpoint, data)
		if err != nil {
			log.Fatalf("Post did not work: %v", err)
		}

		respBody, err := ioutil.ReadAll(req.Body)

		if err != nil {
			log.Fatalf("Could not parse data %v \n", err)
		}

		defer req.Body.Close()

		return string(respBody)
	}()

	fmt.Println("Admin user has been created...")
	return scriptExec, nil
}
