package cmd

import (
	"testing"
	"time"

	client2 "github.com/kubideh/kubesearch/cmd/kubectl-search/client"
	app2 "github.com/kubideh/kubesearch/cmd/kubesearch/app"
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

func createServer(bindAddress string) app2.App {
	appFlags := app2.CreateImmutableServerFlagsWithBindAddress(bindAddress)
	k8sClient := fake.NewSimpleClientset()
	aController := controller.Create(k8sClient)
	anApp := app2.Create(appFlags, aController)
	return anApp
}

func startServer(t *testing.T, anApp app2.App) {
	go func() {
		err := anApp.Run()
		require.NoError(t, err)
	}()
}

func createClient(server string) client2.Client {
	clientFlags := client2.CreateImmutableClientFlagsWithServerAddress(server)
	return client2.Create(clientFlags)
}
