package app

import (
	"net/http"

	"github.com/kubideh/kubesearch/search/finder"
	"github.com/kubideh/kubesearch/search/searcher"
	"github.com/kubideh/kubesearch/search/tokenizer"

	"github.com/kubideh/kubesearch/search/api"
	"github.com/kubideh/kubesearch/search/controller"
	"k8s.io/klog/v2"
)

// ConfigureDefault configures and returns a new App.
func ConfigureDefault() App {
	flags := CreateImmutableServerFlags()
	flags.Parse()

	client := createKubernetesClientset(flags)

	aController := controller.Create(client)

	return Create(flags, aController)
}

// Create returns server App objects.
func Create(flags ImmutableServerFlags, aController *controller.Controller) App {
	aTokenizer := tokenizer.Tokenizer()
	aSearcher := searcher.Create(aController.Index(), aTokenizer)
	aFinder := finder.Create(aController.Store())
	aHandler := api.CreateSearchHandler(aSearcher, aFinder)
	aMux := http.NewServeMux()

	return App{
		controller: aController,
		flags:      flags,
		handler:    aHandler,
		mux:        aMux,
	}
}

// App provides everything needed to run KubeSearch.
type App struct {
	controller *controller.Controller
	flags      ImmutableServerFlags
	handler    http.HandlerFunc
	mux        *http.ServeMux
}

// Run starts the given Controller and registers the Search API
// handler.
func (a App) Run() error {
	// create the Controller to be used by the search API handler
	cancel := a.controller.Start()
	defer cancel()

	api.RegisterSearchHandler(a.mux, a.handler)

	klog.Infoln("Listening on " + a.flags.BindAddress())
	return http.ListenAndServe(a.flags.BindAddress(), a.mux)
}
