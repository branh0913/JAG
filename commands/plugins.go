package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Plugins interface {
	New(endpoint, resource string, plugins map[string]string) (PluginInsecure, error)
	Create() (string, error)
}

type PluginInsecure struct {
	Endpoint string
	Resource string
	Plugins  map[string]string
}

func (p *PluginInsecure) New(endpoint, resource string, plugins map[string]string) (PluginInsecure, error) {

	return PluginInsecure{
		Endpoint: endpoint,
		Resource: resource,
		Plugins:  plugins,
	}, nil
}

func (p PluginInsecure) Create() (string, error) {
	scriptFile := p.Plugins["script"]

	scriptExec := func() string {
		f, err := ioutil.ReadFile(scriptFile)

		if err != nil {
			log.Fatalf("Could not plugin script file %s \n", err)
		}
		data := url.Values{}
		data.Set(p.Resource, string(f))
		req, err := http.PostForm(p.Endpoint, data)
		fmt.Println(req)
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

	return scriptExec, nil

}
