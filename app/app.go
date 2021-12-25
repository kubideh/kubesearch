package app

import (
	"net/http"

	"github.com/kubideh/kubesearch/search/searcher"
	"github.com/kubideh/kubesearch/search/tokenizer"

	"github.com/kubideh/kubesearch/search/api"
	"github.com/kubideh/kubesearch/search/controller"
	"k8s.io/klog/v2"
)

// Create configures and returns a new App.
func Create() App {
	flags := NewFlags()
	flags.Parse()

	client := kubeClientFrom(flags)

	aController := controller.New(client)

	return New(flags, aController)
}

// New returns server App objects.
func New(flags ImmutableFlags, aController *controller.Controller) App {
	return App{
		controller: aController,
		flags:      flags,
		mux:        http.NewServeMux(),
	}
}

// App provides everything needed to run KubeSearch.
type App struct {
	controller *controller.Controller
	flags      ImmutableFlags
	mux        *http.ServeMux
}

// Run starts the given Controller and registers the Search API
// handler.
func (a App) Run() error {
	// create the Controller to be used by the search API handler
	cancel := a.controller.Start()
	defer cancel()

	api.RegisterHandler(a.mux, searcher.Searcher(a.controller.Index(), tokenizer.Tokenizer()), a.controller.Store())

	klog.Infoln("Listening on " + a.flags.BindAddress())
	return http.ListenAndServe(a.flags.BindAddress(), a.mux)
}
