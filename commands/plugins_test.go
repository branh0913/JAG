package commands_test

import (
	"JAG/commands"
	"testing"
	"log"
)


var plugins commands.PluginInsecure
var pluginhash =  map[string]string{"script": "test/somefile"}

func TestPluginNew(t *testing.T)  {



	admininst, err :=  plugins.New("http://example.com", "blah", pluginhash )

	if err != nil{
		t.Fatalf("Could not instantiate admin instance %v %v\n", err, admininst)
	}



}

func TestPlugin_Create(t *testing.T) {


	admininst, err := plugins.New("http://example.com", "blah", pluginhash)



	if err != nil{
		t.Fatalf("Could not instantiate admin instance %v \n", err)
	}

	resp, err := admininst.Create()

	if err != nil{
		log.Fatalf("Could not do request %v  with value: %v \n", err, resp)
	}




}
