// Package api provides the API for searching for Kubernetes
// objects.
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubideh/kubesearch/search"
)

func Search(endpoint, query string) ([]search.Result, error) {
	response, err := http.Get(fmt.Sprintf("%s/v1/search?query=%s", endpoint, query))

	if err != nil {
		return []search.Result{}, err
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	if err != nil {
		return []search.Result{}, err
	}

	var result []search.Result
	if err := json.Unmarshal(body, &result); err != nil {
		return []search.Result{}, err
	}

	return result, nil
}
