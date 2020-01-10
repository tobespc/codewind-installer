package project

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/eclipse/codewind-installer/pkg/connections"
	"github.com/eclipse/codewind-installer/pkg/sechttp"
	"github.com/eclipse/codewind-installer/pkg/utils"
)

type (
	// Project : Represents a project
	Project struct {
		ProjectID      string `json:"projectID"`
		Name           string `json:"name"`
		Language       string `json:"language"`
		Host           string `json:"host"`
		LocationOnDisk string `json:"locOnDisk"`
		AppStatus      string `json:"appStatus"`
	}
)

// GetProject : Get project details from Codewind
func GetProject(httpClient utils.HTTPClient, connection *connections.Connection, url, projectID string) (*Project, error) {
	req, getProjectErr := http.NewRequest("GET", url+"/api/v1/projects/"+projectID+"/", nil)
	if getProjectErr != nil {
		return nil, getProjectErr
	}

	// send request
	resp, httpSecError := sechttp.DispatchHTTPRequest(httpClient, req, connection)
	if httpSecError != nil {
		return nil, httpSecError
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		respErr := errors.New(textAPINotFound)
		return nil, &ProjectError{errOpNotFound, respErr, textAPINotFound}
	}

	byteArray, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	var project Project
	getProjectErr = json.Unmarshal(byteArray, &project)
	if getProjectErr != nil {
		return nil, getProjectErr
	}
	return &project, nil
}

// GetProjectIDFromName : Get a project ID using its name
func GetProjectIDFromName(httpClient utils.HTTPClient, connection *connections.Connection, url, projectName string) (string, error) {
	projects, err := GetAll(httpClient, connection, url)
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.ProjectID, nil
		}
	}
	respErr := errors.New(textAPINotFound)
	return "", &ProjectError{errOpNotFound, respErr, textAPINotFound}
}
