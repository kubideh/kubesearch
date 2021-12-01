// Package api provides the API for searching for Kubernetes
// objects.
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Search is the API used to query for Kubernetes objects.
func Search(endpoint, query string) ([]Result, error) {
	response, err := http.Get(fmt.Sprintf("%s/v1/search?query=%s", endpoint, query))

	if err != nil {
		return []Result{}, err
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	if err != nil {
		return []Result{}, err
	}

	var result []Result
	if err := json.Unmarshal(body, &result); err != nil {
		return []Result{}, err
	}

	return result, nil
}
