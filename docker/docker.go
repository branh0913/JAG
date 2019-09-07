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

type DockerBuildRequest struct {
	Name        string
	ArchiveFile string
	DockerFile  string
	Tag         string
	tarfile     *archivex.TarFile
}

type BuildResults struct {
	Id  string
	Tag string
}

func NewDockerBuild(config *DockerBuildRequest) (*DockerBuildRequest, error) {

	tar := archivex.TarFile{}
	return &DockerBuildRequest{
		Name:        config.Name,
		ArchiveFile: config.ArchiveFile,
		DockerFile:  config.DockerFile,
		tarfile:     &tar,
	}, nil

}

func (dr *DockerBuildRequest) BuildDocker() (*BuildResults, error) {

	buildResults := BuildResults{}

	var tarFile *os.File
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Println("could not create client ", err)
	}

	archiveQuery, ok := archiveExist(dr.ArchiveFile)
	tarFile = archiveQuery
	if !ok {
		tar, err := createArchive(dr.ArchiveFile, dr.tarfile)
		if err != nil {
			return nil, fmt.Errorf("could not create new archive %v", err)
		}
		if err := addDocker(dr.DockerFile, dr.tarfile);err != nil{
			return nil, fmt.Errorf("error when adding docker file to archive %v", err)
		}

		tarFile = tar

	}

	imgBuildResp, err := dockerClient.ImageBuild(context.Background(), tarFile, types.ImageBuildOptions{
		Tags: []string{dr.Tag},
	})

	if err != nil {
		return nil, fmt.Errorf("could not build %v", err)
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

func archiveExist(tarFile string) (*os.File, bool) {
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		return nil, false
	}

	tar, err := os.Open(tarFile)
	if err != nil {
		return nil, true
	}

	return tar, true

}

func createArchive(tarFile string, tararch *archivex.TarFile) (*os.File, error) {

	if err := tararch.Create(tarFile); err != nil {
		return nil, fmt.Errorf("error occured when creating archive %v", err)
	}

	defer tararch.Close()

	tar, err := os.Open(tarFile)
	if err != nil {
		return nil, fmt.Errorf("could not create new file %v", err)
	}

	if err := tar.Close(); err != nil {
		return nil, fmt.Errorf("file not able to close %v", err)
	}

	return tar, nil

}

func addDocker(dockerFile string, tarArch *archivex.TarFile) error {
	f, err := os.Open(dockerFile)
	if err != nil {
		return nil
	}

	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info %v", err)
	}

	if err := tarArch.Add("Dockerfile", f, fInfo); err != nil{
		return fmt.Errorf("could not add dockerfile to archive %v", err)
	}

	return nil

}
