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

// CreateDefault configures and returns a new App.
func CreateDefault() App {
	flags := CreateImmutableServerFlags()
	flags.Parse()

	client := createKubernetesClientset(flags)

	aController := controller.Create(client)

	return CreateFromFlagsAndController(flags, aController)
}

// CreateFromFlagsAndController returns server App objects.
func CreateFromFlagsAndController(flags ImmutableServerFlags, aController *controller.Controller) App {
	return App{
		controller: aController,
		flags:      flags,
		mux:        http.NewServeMux(),
	}
}

// App provides everything needed to run KubeSearch.
type App struct {
	controller *controller.Controller
	flags      ImmutableServerFlags
	mux        *http.ServeMux
}

// Run starts the given Controller and registers the Search API
// handler.
func (a App) Run() error {
	// create the Controller to be used by the search API handler
	cancel := a.controller.Start()
	defer cancel()

	api.RegisterHandler(a.mux, searcher.Create(a.controller.Index(), tokenizer.Tokenizer()), finder.Create(a.controller.Store()))

	klog.Infoln("Listening on " + a.flags.BindAddress())
	return http.ListenAndServe(a.flags.BindAddress(), a.mux)
}
