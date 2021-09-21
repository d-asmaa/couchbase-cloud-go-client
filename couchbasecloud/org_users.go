package couchbasecloud

import (
	"net/http"
	"net/url"
	"strconv"
)

type Users []Users

type UsersList struct {
	Cursor Cursor  `json:"cursor"`
	Data   []Users `json:"data"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type ListUsersOptions struct {
	Page    int     `json:"page"`
	PerPage int     `json:"perPage"`
	SortBy  *string `json:"sortBy"`
}

const usersUrl = "/users"

func (client *CouchbaseCloudClient) ListUsers(options *ListUsersOptions) (*UsersList, error) {
	cloudsUrl := client.BaseURL + client.getApiEndpoint(usersUrl)

	if options != nil {
		setListUsersParams(&cloudsUrl, *options)
	}

	req, err := http.NewRequest(http.MethodGet, cloudsUrl, nil)
	if err != nil {
		return nil, err
	}

	res := UsersList{}

	if err := client.sendRequest(req, &res, true); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListUserPages allows iterating over all the users. For every page of user items it will call the callback and pass
// the page worth of users as well as a boolean that indicates whether is is the last page or not.
// The function iterates over all the pages either until the callback returns false, the REST endpoint returns an error
// or it runs out of pages.
func (client *CouchbaseCloudClient) ListUserPages(options *ListUsersOptions, fn func(Users, bool) bool) error {
	var localOpts ListUsersOptions
	if options != nil {
		localOpts = *options
	}

	for {
		users, err := client.ListUsers(&localOpts)
		if err != nil {
			return err
		}

		if len(users.Data) == 0 {
			return nil
		}

		cont := fn(users.Data, users.Cursor.Pages.Last >= users.Cursor.Pages.Page)
		if !cont {
			return nil
		}

		localOpts.Page++
	}
}

func setListUsersParams(urlStr *string, options ListUsersOptions) {
	params := url.Values{}

	if options.SortBy != nil {
		params.Add("sortBy", *options.SortBy)
	}

	if options.Page != 0 {
		params.Add("page", strconv.Itoa(options.Page))
	}

	if options.PerPage != 0 {
		params.Add("perPage", strconv.Itoa(options.PerPage))
	}

	if urlParams := params.Encode(); urlParams != "" {
		*urlStr += "?" + urlParams
	}
}
