// gitreader.go
package utils

import (
	"archive/zip"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GitReader struct {
	Name        string
	Url         string
	Version     string
	ProjectPath string
}

func NewGitReader(projectName string, record DependencyRecord) *GitReader {
	var gitReader GitReader

	gitReader.Name = projectName
	gitReader.Url = record.Location
	gitReader.Version = record.Version

	gitReader.ProjectPath = record.ProjectPath

	return &gitReader

}

func (reader *GitReader) CloneDependency() error {
	zipFileName := reader.Name + "." + reader.Version + ".zip"
	zipFilePath := reader.Url + "/zipball" + "/" + reader.Version

	os.MkdirAll(reader.ProjectPath, os.ModePerm)
	output, err := os.Create(path.Join(reader.ProjectPath, zipFileName))
	Check(err)
	defer output.Close()

	response, err := http.Get(zipFilePath)
	Check(err)
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	Check(err)

	zipReader, err := zip.OpenReader(path.Join(reader.ProjectPath, zipFileName))
	Check(err)

	u, err := url.Parse(reader.Url)
	team, project := path.Split(u.Path)

	team = strings.TrimPrefix(team, "/")
	team = strings.TrimSuffix(team, "/")

	for _, file := range zipReader.File {
		path := filepath.Join(reader.ProjectPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		Check(err)
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		Check(err)
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	rootDir := filepath.Join(reader.ProjectPath, zipReader.File[0].Name)
	depPath := filepath.Join(reader.ProjectPath, project+"-"+reader.Version)
	os.Rename(rootDir, depPath)
	os.Remove(filepath.Join(reader.ProjectPath, zipFileName))

	return nil
}
