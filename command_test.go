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

	anApp := createServer(":" + port)
	startServer(t, anApp)

	aClient := createClient("localhost:" + port)

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

func createServer(bindAddress string) app.App {
	appFlags := app.CreateImmutableServerFlagsWithBindAddress(bindAddress)
	k8sClient := fake.NewSimpleClientset()
	aController := controller.Create(k8sClient)
	anApp := app.Create(appFlags, aController)
	return anApp
}

func startServer(t *testing.T, anApp app.App) {
	go func() {
		err := anApp.Run()
		require.NoError(t, err)
	}()
}

func createClient(server string) client.Client {
	clientFlags := client.CreateImmutableClientFlagsWithServerAddress(server)
	return client.Create(clientFlags)
}
