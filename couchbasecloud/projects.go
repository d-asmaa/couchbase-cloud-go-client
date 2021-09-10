package couchbasecloud

import (
	"net/http"
	"net/url"
	"strconv"
)

type Projects []Project

type ProjectsList struct {
	Cursor Cursor    `json:"cursor"`
	Data   []Project `json:"data"`
}

type Project struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	TenantId  string `json:"tenantId"`
	CreatedAt string `json:"createdAt"`
}

type ListProjectsOptions struct {
	Page    int     `json:"page"`
	PerPage int     `json:"perPage"`
	SortBy  *string `json:"sortBy"`
}

type CreateProjectPayload struct {
	Name string `json:"name"`
}

const projectsUrl = "/projects"

func (client *CouchbaseCloudClient) ListProjects(options *ListProjectsOptions) (*ProjectsList, error) {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(projectsUrl)

	if options != nil {
		setListProjectsParams(&cloudsUrl, *options)
	}

	req, err := http.NewRequest(http.MethodGet, cloudsUrl, nil)
	if err != nil {
		return nil, err
	}

	res := ProjectsList{}

	if err := client.sendRequest(req, &res, true); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListProjectPages allows iterating over all the projects. For every page of project items it will call the callback
// and pass the page worth of clouds as well as a boolean that indicates whether is is the last page or not. The
// function iterates over all the pages either until the callback returns false, the REST endpoint returns an error
// or it runs out of pages.
func (client *CouchbaseCloudClient) ListProjectPages(opts *ListProjectsOptions, fn func(Projects, bool) bool) error {
	var localOpts ListProjectsOptions
	if opts != nil {
		localOpts = *opts
	}

	for {
		clouds, err := client.ListProjects(&localOpts)
		if err != nil {
			return err
		}

		if len(clouds.Data) == 0 {
			return nil
		}

		cont := fn(clouds.Data, clouds.Cursor.Pages.Last >= clouds.Cursor.Pages.Page)
		if !cont {
			return nil
		}

		localOpts.Page++
	}
}

func setListProjectsParams(projectsUrl *string, options ListProjectsOptions) {
	params := url.Values{}

	if options.Page != 0 {
		params.Add("page", strconv.Itoa(options.Page))
	}

	if options.PerPage != 0 {
		params.Add("perPage", strconv.Itoa(options.PerPage))
	}

	if options.SortBy != nil {
		params.Add("sortBy", *options.SortBy)
	}

	if urlParams := params.Encode(); urlParams != "" {
		*projectsUrl += "?" + urlParams
	}
}

func (client *CouchbaseCloudClient) CreateProject(payload *CreateProjectPayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(projectsUrl)

	req, err := http.NewRequest(http.MethodPost, cloudsUrl, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, payload, true); err != nil {
		return err
	}

	return nil
}

func (client *CouchbaseCloudClient) DeleteProject(project *Project) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(projectsUrl)

	req, err := http.NewRequest(http.MethodDelete, cloudsUrl+"/"+project.Id, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, nil, false); err != nil {
		return err
	}

	return nil
}
