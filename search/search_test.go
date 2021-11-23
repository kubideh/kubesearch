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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*httptest.Server, context.CancelFunc) {
	index := make(map[string]string)
	SetIndex(index)

	podQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-queue")
	go func(podQueue workqueue.RateLimitingInterface) {
		key, shutdown := podQueue.Get()
		for !shutdown {
			var name string

			metadata := strings.Split(key.(string), "/")
			if len(metadata) == 1 {
				name = metadata[0]
			} else {
				name = metadata[1]
			}

			index[name] = key.(string)

			key, shutdown = podQueue.Get()
		}
		klog.Info("Shutting down pod-queue")
	}(podQueue)

	client := fake.NewSimpleClientset()

	informerFactory := informers.NewSharedInformerFactory(client, 0)

	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Error(err)
			} else {
				podQueue.Add(key)
			}
		},
	})
	SetInformer(podInformer)

	ctx, cancel := context.WithCancel(context.Background())
	informerFactory.Start(ctx.Done())

	_, err := client.CoreV1().Pods("flargle").Create(
		context.TODO(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "blargle",
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	server := httptest.NewServer(http.DefaultServeMux)

	return server, cancel
}

func TestSearch_podByName(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	params := url.Values{}
	params.Add("query", "blargle")

	response, err := http.Get(fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode()))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{"kind":"Pods","namespace":"flargle","name":"blargle"}`, string(body))
}

func TestSearch_nonExistentPod(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	params := url.Values{}
	params.Add("query", "whatever")

	response, err := http.Get(fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode()))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{}`, string(body))
}

func TestSearch_missingQuery(t *testing.T) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	response, err := http.Get(fmt.Sprintf("%s/v1/search", server.URL))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{}`, string(body))
}
