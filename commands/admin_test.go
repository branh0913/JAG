package commands_test

import (
	"testing"
	"JAG/commands"
	"log"
)

var admin commands.Admin

func TestAdminNew(t *testing.T)  {



	admininst, err := admin.New("http://example.com", "blah", "john", "doe")

	if err != nil{
		t.Fatalf("Could not instantiate admin instance %v with value %v\n", err, admininst)
	}

}

func TestAdmin_Create(t *testing.T) {
	admininst, err := admin.New("http://example.com", "blah", "john", "doe")

	if err != nil{
		t.Fatalf("Could not instantiate admin instance %v \n", err)
	}

	resp, err := admininst.Create()

	if err != nil{
		log.Fatalf("Could not do request %v with value %v\n", err, resp)
	}




}


