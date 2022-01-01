// Package api provides the API for searching for Kubernetes
// objects.
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Search is the API used to queryString for Kubernetes objects.
func Search(endpoint, query string) (result []Result, err error) {
	response, err := http.Get(searchURL(endpoint, query))

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)

	return
}

func searchURL(endpoint, query string) string {
	return fmt.Sprintf("%s%s?%s=%s", endpoint, endpointPath, queryParamName, url.QueryEscape(query))
}
