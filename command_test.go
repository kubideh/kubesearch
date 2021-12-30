package main

import (
	"testing"
	"time"

	"github.com/kubideh/kubesearch/app"
	"github.com/kubideh/kubesearch/client"
	"github.com/kubideh/kubesearch/search/controller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCommand(t *testing.T) {
	const port = "31337"

	configureServer(t, ":"+port)
	aClient := configureClient("localhost:" + port)

	var clientErr error
	for i := 1; i <= 3; i++ {
		clientErr = aClient.Run()

		if clientErr == nil {
			break
		}

		time.Sleep(time.Duration(i) * time.Second)
	}

	assert.NoError(t, clientErr)
}

func configureServer(t *testing.T, bindAddress string) {
	appFlags := app.NewFlagsWithBindAddress(bindAddress)
	k8sClient := fake.NewSimpleClientset()
	aController := controller.New(k8sClient)
	anApp := app.New(appFlags, aController)

	go func() {
		err := anApp.Run()
		require.NoError(t, err)
	}()
}

func configureClient(server string) client.Client {
	clientFlags := client.NewFlagsWithServer(server)
	return client.New(clientFlags)
}
