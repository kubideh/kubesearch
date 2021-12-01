package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubideh/kubesearch/search/controller"
	"github.com/kubideh/kubesearch/search/index"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type testSearchCase struct {
	name   string
	query  string
	result []Result
}

func TestSearchAPI(t *testing.T) {
	cases := []testSearchCase{
		{
			name:  "search for pod by name",
			query: "blargle",
			result: []Result{
				{
					Kind:      "Pod",
					Name:      "blargle",
					Namespace: "flargle",
				},
			},
		},
		{
			name:  "search for pod by namespace",
			query: "flargle",
			result: []Result{
				{
					Kind:      "Pod",
					Name:      "blargle",
					Namespace: "flargle",
				},
				{
					Kind:      "Pod",
					Name:      "foo",
					Namespace: "flargle",
				},
			},
		},
		{
			name:   "search for missing object",
			query:  "whatever",
			result: []Result{},
		},
		{
			name:   "search using empty query",
			query:  "",
			result: []Result{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { testSearch(t, c) })
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

	cont := controller.NewController(client)

	idx := index.NewIndex()

	cancel := cont.Start(idx)
	defer cancel()

	mux := http.NewServeMux()
	RegisterHandler(mux, idx, cont.Store())

	server := httptest.NewServer(mux)
	defer server.Close()

	for _, p := range testPods() {
		_, err := client.CoreV1().Pods(p.GetNamespace()).Create(context.TODO(), p, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	t.Log("query: ", c.query)

	result, err := Search(server.URL, c.query)

	t.Log("result: ", c.result)

	assert.NoError(t, err)
	assert.Equal(t, c.result, result)
}
