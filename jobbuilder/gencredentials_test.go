package jobbuilder_test

import (
	"testing"
	"JAG/jobbuilder"
)

func TestCredentialID_Create(t *testing.T) {

	test := jobbuilder.CredentialID{}
	var ptrstr *string
	sa_user := "ZmFrZQ=="
	sa_pass := "ZmFrZXBhc3M="
	testNew, err:= test.New("http://example.com", "blah", "blah", "pass",sa_user,sa_pass,
		                     ptrstr)

	if err != nil{
		t.Errorf("Could not instantiate %v %v", err, testNew)
	}
}

func TestCredentialID_Create2(t *testing.T) {
	test := jobbuilder.CredentialID{}
	var ptrstr *string
	sa_user := "ZmFrZQ=="
	sa_pass := "ZmFrZXBhc3M="
	testNew, err:= test.New("http://example.com", "blah", "blah", "pass",sa_user,sa_pass,
		ptrstr)

	if err != nil{
		t.Errorf("Could not instantiate %v %v", err, testNew)
	}

	testcreate, err := testNew.Create()

	if err != nil{
		t.Errorf("Could not generate credentialid %v %v", err, testcreate)
	}
}

