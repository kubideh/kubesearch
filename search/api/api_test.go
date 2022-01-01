package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubideh/kubesearch/search/finder"
	"github.com/kubideh/kubesearch/search/searcher"
	"github.com/kubideh/kubesearch/search/tokenizer"

	"github.com/kubideh/kubesearch/search/controller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSearch_emptyQuery(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	result, err := Search(server.URL, "")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSearch_queryForMissingObjects(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	result, err := Search(server.URL, "whatever")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSearch_queryUsingMultipleTerms(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	result, err := Search(server.URL, "search for something")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSearch_queryForSinglePod(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	result, err := Search(server.URL, "blargle")

	assert.NoError(t, err)
	assert.Equal(t, []Result{
		{
			Kind:      "Pod",
			Name:      "blargle",
			Namespace: "flargle",
			Rank:      1,
		},
	}, result)
}

func TestSearch_queryForAllPodsInNamespace(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	result, err := Search(server.URL, "flargle")

	assert.NoError(t, err)
	assert.Equal(t, []Result{
		{
			Kind:      "Pod",
			Name:      "blargle",
			Namespace: "flargle",
			Rank:      1,
		},
		{
			Kind:      "Pod",
			Name:      "foo",
			Namespace: "flargle",
			Rank:      1,
		},
	}, result)
}

func setup(t *testing.T) (*httptest.Server, context.CancelFunc) {
	client := fake.NewSimpleClientset()

	aController := controller.Create(client)
	cancel := aController.Start()

	objectFinder := finder.Create(aController.Store())

	mux := http.NewServeMux()
	RegisterHandler(mux, searcher.Create(aController.Index(), tokenizer.Tokenizer()), objectFinder)

	server := httptest.NewServer(mux)

	for _, p := range testPods() {
		_, err := client.CoreV1().Pods(p.GetNamespace()).Create(context.TODO(), p, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	return server, cancel
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
