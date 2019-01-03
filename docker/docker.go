package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type BuildResults struct {
	Id  string
	Tag string
}

func BuildDocker() (*BuildResults, error) {

	buildResults := BuildResults{}

	var tarFile *os.File
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Println("could not create client ", err)
	}

	if _, err := os.Stat("tmp/myfile.tar"); os.IsNotExist(err) {
		tararch := archivex.TarFile{}
		f, err := os.Open("Dockerfile")
		if err != nil {
			logrus.Fatalf("could not open dockerfile %v", err)
		}
		fInfo, err := f.Stat()
		if err != nil {
			logrus.Fatalf("could not get file info %v", err)
		}
		tararch.Create("tmp/myfile.tar")
		tararch.Add("Dockerfile", f, fInfo)
		tararch.Close()
		tarFile, err = os.Open("tmp/myfile.tar")
		if err != nil {
			fmt.Println(err)
			logrus.Fatalf("could not create file %v", err)
		}

		defer tarFile.Close()

	}

	imgBuildResp, err := dockerClient.ImageBuild(context.Background(), tarFile, types.ImageBuildOptions{
		Tags: []string{"jenkinstest"},
	})

	if err != nil {
		logrus.Errorf("could not build %v", err)
	}

	defer imgBuildResp.Body.Close()

	_, err = ioutil.ReadAll(imgBuildResp.Body)
	if err != nil {
		logrus.Errorf("could not read bytes due to %v", err)
	}

	listImages, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{
		All:     false,
		Filters: filters.Args{},
	})
	if err != nil {
		logrus.Errorf("could not list images due to %v", err)
	}

	for _, images := range listImages {
		if strings.Contains(images.RepoTags[0], "jenkinstest") {
			buildResults.Id = images.ID
			buildResults.Tag = images.RepoTags[0]

		}
	}

	return &buildResults, nil
}
