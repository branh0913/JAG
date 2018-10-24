package jobbuilder_test

import (
	"github.com/brharrelldev/jag/jobbuilder"
	"os"
	"testing"
)

type MockedFileInterface interface {
	Create(name string) (*MockFile, error)
}

type MockedJbuilder struct {
	jobbuilder.JBuilderConfig
	endpoint string
	apitoken string
	jconfig  string
}

type MockFile struct {
	os.File
}

func (m MockedJbuilder) Create(name string) (*MockFile, error) {

	return &MockFile{}, nil
}

func TestJBuilderConfig_New(t *testing.T) {

	config := new(jobbuilder.JBuilderConfig)

	jb, err := config.New("http://example.com", "11111111111", "/tmp/test.ini")

	if err != nil {
		t.Fatalf("Could not instantiate object %v with value %v", jb, err)
	}
}

func TestJBuilderConfig_BuildFile(t *testing.T) {

	newOS := &MockFile{}
	oldOS := os.File{}
	oldOS = newOS.File

	defer func() {
		oldOS = os.File{}
	}()

	jb := jobbuilder.JBuilderConfig{}
	config, err := jb.New("http://example.com", "11111111", "/tmp/test.ini")

	if err != nil {
		t.Errorf("Could not be intantiated %v \n", err)
	}

	config.BuildFile("johndoe", "templates/jenkins.ini.gotmpl")
}
