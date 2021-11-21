package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/kubideh/kubesearch/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_usingInformer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kubesearch")
	defer queue.ShutDown()

	aClient := fake.NewSimpleClientset()

	informerFactory := informers.NewSharedInformerFactory(aClient, 0)

	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			} else {
				t.Log(err)
			}
			pod := obj.(*corev1.Pod)
			t.Logf("pod added: %s/%s", pod.Namespace, pod.Name)
		},
	})

	informerFactory.Start(ctx.Done())

	cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	_, err := aClient.CoreV1().Pods("flargle").Create(
		context.TODO(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "blargle",
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	key, _ := queue.Get()

	queue.Done(key)

	item, exists, err := podInformer.GetStore().GetByKey(key.(string))
	if !exists || err != nil {
		queue.Forget(key)
		require.NoError(t, err)
	}

	pod := item.(*corev1.Pod)
	assert.Equal(t, "flargle", pod.Namespace)
	assert.Equal(t, "blargle", pod.Name)
}

func TestSearch_podByName(t *testing.T) {
	aClient := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "flargle",
			Name:      "blargle",
		},
	})
	client.SetClient(aClient)

	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

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
	aClient := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "flargle",
			Name:      "blargle",
		},
	})
	client.SetClient(aClient)

	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

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
	aClient := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "flargle",
			Name:      "blargle",
		},
	})
	client.SetClient(aClient)

	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	response, err := http.Get(fmt.Sprintf("%s/v1/search", server.URL))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{}`, string(body))
}
