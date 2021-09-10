package couchbasecloud

import (
	"net/http"
)

type BucketsList struct {
	Cursor Cursor   `json:"cursor"`
	Data   []Bucket `json:"data"`
}

type Bucket struct {
	Name               string `json:"name"`
	MemoryQuota        int    `json:"memoryQuota"`
	Replicas           int    `json:"replicas"`
	ConflictResolution string `json:"conflictResolution"` // TODO: enum type
	Status             string `json:"string"`
}

type BucketCreatePayload struct {
	Name        string `json:"name"`
	MemoryQuota int    `json:"memoryQuota"`
}

type BucketDeletePayload struct {
	Name string `json:"name"`
}

func (client *CouchbaseCloudClient) ListBuckets(cluster *Cluster) (*BucketsList, error) {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/buckets")

	req, err := http.NewRequest(http.MethodGet, cloudsUrl, nil)
	if err != nil {
		return nil, err
	}

	res := BucketsList{}

	if err := client.sendRequest(req, &res, true); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *CouchbaseCloudClient) CreateBucket(cluster *Cluster, payload *BucketCreatePayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/buckets")

	req, err := http.NewRequest(http.MethodPost, cloudsUrl, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, payload, true); err != nil {
		return err
	}

	return nil
}

func (client *CouchbaseCloudClient) DeleteBucket(cluster *Cluster, payload *BucketDeletePayload) error {
	cloudsUrl := client.BaseURL + client.getApiEndpoint("/clusters/"+cluster.Id+"/buckets")

	req, err := http.NewRequest(http.MethodDelete, cloudsUrl, nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, payload, true); err != nil {
		return err
	}

	return nil
}
