// Package api provides the API for searching for Kubernetes
// objects.
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Search(endpoint, query string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%s/v1/search?query=%s", endpoint, query))

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	if err != nil {
		return "", err
	}

	return string(body), nil
}
