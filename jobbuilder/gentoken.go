package jobbuilder

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type APIToken interface {
	New(endpoint, resource, username, password string) (*Token, error)
	Retrieve(user string) string
}

type Token struct {
	Endpoint string
	Resource string
	Username string
	Password string
}

func (t *Token) New(endpoint, resource, username, password string) (*Token, error) {

	return &Token{
		Endpoint: endpoint,
		Resource: resource,
		Username: username,
		Password: password,
	}, nil
}

func (t Token) Retrieve(user string) string {

	type Currentuser struct {
		Currentuser string
	}

	userobj := Currentuser{Currentuser: user}

	groovyScript := `user = hudson.model.User.get("{{.Currentuser}}")
prop = user.getProperty(jenkins.security.ApiTokenProperty.class)
println(prop.getApiTokenInsecure())`

	tmpl := template.Must(template.New("APIToken").Parse(groovyScript))

	var apitoken bytes.Buffer

	tmplinvoke := tmpl.Execute(&apitoken, userobj)

	if tmplinvoke != nil {
		log.Fatal("Error when trying to invoke user lookup template %v \n", tmplinvoke)
	}

	log.Printf("Starting retrieval of API Token...")

	tokenresp := func() string {

		client := &http.Client{}
		u := url.Values{}
		u.Set(t.Resource, apitoken.String())

		req, err := http.NewRequest(http.MethodPost, t.Endpoint, strings.NewReader(u.Encode()))

		if err != nil {
			log.Fatalf("Request object could not be build %v", err)
		}

		req.SetBasicAuth(t.Username, t.Password)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		clientresp, err := client.Do(req)

		if err != nil {
			log.Fatalf("Request failed with %v \n", err)
		}

		defer clientresp.Body.Close()
		respString, err := ioutil.ReadAll(clientresp.Body)

		return string(respString)
	}()

	return tokenresp

}
