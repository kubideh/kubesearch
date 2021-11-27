package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSearchCase struct {
	name   string
	params string
	result string
}

func TestSearch(t *testing.T) {
	cases := []testSearchCase{
		{
			name:   "search for objects by name",
			params: "query=blargle",
			result: `[{"kind":"Deployment","namespace":"flargle","name":"blargle"},{"kind":"Pod","namespace":"flargle","name":"blargle"}]`,
		},
		{
			name:   "search for objects by namespace",
			params: "query=flargle",
			result: `[{"kind":"Deployment","namespace":"flargle","name":"blargle"},{"kind":"Pod","namespace":"flargle","name":"blargle"},{"kind":"Pod","namespace":"flargle","name":"foo"}]`,
		},
		{
			name:   "search for missing object",
			params: "query=whatever",
			result: `[]`,
		},
		{
			name:   "search using empty query",
			params: "query=",
			result: `[]`,
		},
		{
			name:   "search with missing query param",
			params: "",
			result: `[]`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { testSearch(t, c) })
	}
}

func testDeployments() []*appsv1.Deployment {
	return []*appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "blargle",
				Namespace: "flargle",
			},
		},
	}
}

func testPods() []*corev1.Pod {
	return []*corev1.Pod{
		testPodFlargleBlargle(),
		testPodFlargleFoo(),
	}
}

func testPodFlargleBlargle() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "blargle",
			Namespace: "flargle",
		},
	}
}

func testPodFlargleFoo() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "flargle",
		},
	}
}

func testSearch(t *testing.T, c testSearchCase) {
	client := fake.NewSimpleClientset()

	controller := NewController(client)

	index := NewIndex()

	cancel := controller.Start(index)
	defer cancel()

	mux := http.NewServeMux()
	RegisterHandler(mux, index, controller.Store())

	server := httptest.NewServer(mux)
	defer server.Close()

	for _, d := range testDeployments() {
		_, err := client.AppsV1().Deployments(d.GetNamespace()).Create(context.TODO(), d, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	for _, p := range testPods() {
		_, err := client.CoreV1().Pods(p.GetNamespace()).Create(context.TODO(), p, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	params := url.Values{}
	if c.params != "" {
		params.Add(strings.Split(c.params, "=")[0], strings.Split(c.params, "=")[1])
	}

	uri := fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode())
	t.Log(uri)
	response, err := http.Get(uri)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, c.result, string(body))
}
