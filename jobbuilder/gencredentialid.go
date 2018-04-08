package jobbuilder

import (
	"html/template"
	"log"
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
)

type Credentials interface {
	New(endpoint, resource, username, password, serviceaccount, servicepassword string,  credentialid *string) (*CredentialID, error)
	Create()(string, error)
}

type CredentialID struct {
	Endpoint string
	Resource string
	Username string
	Password string
	ServiceAccount string
	Servicepass  string
	CredentialID *string
}

func (c *CredentialID) New(endpoint, resource, username, password, serviceaccount, servicepass string, credentialid *string) (*CredentialID, error)  {

	serviceaccountdecode, err := base64.StdEncoding.DecodeString(serviceaccount)


	if err != nil{
		log.Fatalf("Could not decode string %v", err)
	}

	servicepassdecode, err := base64.StdEncoding.DecodeString(servicepass)

	if err != nil{
		log.Fatalf("Could not decode string %v", err)
	}

	return &CredentialID{
		Endpoint: endpoint,
		Resource: resource,
		Username: username,
		Password: password,
		ServiceAccount: string(serviceaccountdecode),
		Servicepass: string(servicepassdecode),
		CredentialID: credentialid,
	}, nil

}

func (c CredentialID) Create()(string, error) {


	groovyScript, err := template.ParseFiles("groovyScripts/createCredentials.groovy")

	if err != nil{
		log.Fatalf("File could not be found or template could not be generated %v", err)
	}

	var credentialid bytes.Buffer

	groovyScript.Execute(&credentialid, c)

	gencid := func () (string){

		client := &http.Client{}
		u := url.Values{}
		u.Set(c.Resource, credentialid.String())

		req,err := http.NewRequest(http.MethodPost, c.Endpoint, strings.NewReader(u.Encode()))

		if err != nil{
			log.Fatalf("Request object could not be build %v", err)
		}

		req.SetBasicAuth(c.Username,c.Password)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		clientresp, err := client.Do(req)


		if err != nil{
			log.Fatalf("Request failed with %v \n", err)
		}

		defer clientresp.Body.Close()
		respString, err := ioutil.ReadAll(clientresp.Body)

		return string(respString)
	}()

	return gencid, nil



}