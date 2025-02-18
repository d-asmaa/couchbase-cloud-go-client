package couchbasecloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Clusters []Cluster

type ClustersList struct {
	Cursor Cursor    `json:"cursor"`
	Data   []Cluster `json:"data"`
}

type Cluster struct {
	Id        string   `json:"id"`
	CloudId   string   `json:"cloudId"`
	Name      string   `json:"name"`
	Nodes     int      `json:"nodes"`
	ProjectId string   `json:"projectId"`
	Services  []string `json:"services"`
	TenantId  string   `json:"tenantId"`
}

type ListClustersOptions struct {
	Page      int     `json:"page"`
	PerPage   int     `json:"perPage"`
	SortBy    *string `json:"sortBy"`
	CloudId   *string `json:"cloudId"`
	ProjectId *string `json:"projectId"`
}

type ClusterCreatePayload struct {
	Name      string `json:"name`
	CloudId   string `json:"cloudId"`
	ProjectId string `json:"projectId"`
	// TO ADD:
	// Servers        []Server                            `json:servers`
	// SupportPackage *CreateClusterRequestSupportPackage `json:supportPackage`
	// Version        *ClusterVersions                    `json:version`
}

const clustersUrl = "/clusters"

func (client *CouchbaseCloudClient) ListClusters(options *ListClustersOptions) (*ClustersList, error) {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(clustersUrl)

	if options != nil {
		setListClustersParams(&cloudsUrl, *options)
	}

	req, err := http.NewRequest(http.MethodGet, cloudsUrl, nil)
	if err != nil {
		return nil, err
	}

	res := ClustersList{}

	if err := client.sendRequest(req, &res, true); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListClusterPages allows iterating over all the clusters. For every page of cluster items it will call the callback
// and pass the page worth of clouds as well as a boolean that indicates whether is is the last page or not. The
// function iterates over all the pages either until the callback returns false, the REST endpoint returns an error
// or it runs out of pages.
func (client *CouchbaseCloudClient) ListClusterPages(opts *ListClustersOptions, fn func(Clusters, bool) bool) error {
	var localOpts ListClustersOptions
	if opts != nil {
		localOpts = *opts
	}

	for {
		clouds, err := client.ListClusters(&localOpts)
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

func setListClustersParams(clustersUrl *string, options ListClustersOptions) {
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

	if options.CloudId != nil {
		params.Add("cloudId", *options.CloudId)
	}

	if options.ProjectId != nil {
		params.Add("projectId", *options.ProjectId)
	}

	if urlParams := params.Encode(); urlParams != "" {
		*clustersUrl += "?" + urlParams
	}
}

func (client *CouchbaseCloudClient) CreateCluster(payload *ClusterCreatePayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(clustersUrl)

	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(http.MethodPost, cloudsUrl, body)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, payload, true); err != nil {
		return err
	}

	return nil
}

func (client *CouchbaseCloudClient) DeleteCluster(cluster *Cluster) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(clustersUrl)

	req, err := http.NewRequest(http.MethodDelete, cloudsUrl+"/"+cluster.Id, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, nil, false); err != nil {
		return err
	}

	return nil
}
