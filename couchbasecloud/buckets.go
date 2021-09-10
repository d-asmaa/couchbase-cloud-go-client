package couchbasecloud

import (
	"net/http"
)

type BucketCreatePayload struct {
	Name        string `json:"name"`
	MemoryQuota int    `json:"memoryQuota"`
}

type BucketDeletePayload struct {
	Name string `json:"name"`
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
