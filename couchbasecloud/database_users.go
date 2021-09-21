package couchbasecloud

import (
	"net/http"
)

type DatabaseUsers []DatabaseUser

type DatabaseUsersList struct {
	Cursor Cursor         `json:"cursor"`
	Data   []DatabaseUser `json:"data"`
}

type DatabaseUser struct {
	UserId   string       `json:userId,omitempty"`
	Username string       `json:"username,omitempty"`
	Password string       `json:"password,omitempty"`
	Access   []BucketRole `json:"buckets,omitempty"`
}

type DatabaseUserCreatePayload struct {
	UserId           string       `json:userId,omitempty"`
	Username         string       `json:"username,omitempty"`
	Password         string       `json:"password,omitempty"`
	Buckets          []BucketRole `json:"buckets,omitempty"`
	AllBucketsAccess string       `json:"allBucketsAccess,omitempty"`
}

type DatabaseUserDeletePayload struct {
	Username string `json:"username,omitempty"`
}

type BucketRole struct {
	Name  string   `json:"bucketName"`
	Roles []string `json:"bucketAccess"`
}

func (client *CouchbaseCloudClient) ListDatabaseUsers(cluster *Cluster) (*DatabaseUsersList, error) {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/users")

	req, err := http.NewRequest(http.MethodGet, cloudsUrl, nil)
	if err != nil {
		return nil, err
	}

	res := DatabaseUsersList{}

	if err := client.sendRequest(req, &res, true); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *CouchbaseCloudClient) CreateDatabaseUser(cluster *Cluster, payload *DatabaseUserCreatePayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/users")

	req, err := http.NewRequest(http.MethodPost, cloudsUrl, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, payload, true); err != nil {
		return err
	}

	return nil
}

func (client *CouchbaseCloudClient) DeleteDatabaseUser(cluster *Cluster, payload *DatabaseUserDeletePayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/users"+payload.Username)

	req, err := http.NewRequest(http.MethodDelete, cloudsUrl, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, nil, false); err != nil {
		return err
	}

	return nil
}
