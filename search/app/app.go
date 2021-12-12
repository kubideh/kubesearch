package app

import (
	"flag"
	"github.com/kubideh/kubesearch/search/index"
	"net/http"
	"path/filepath"

	"github.com/kubideh/kubesearch/search/api"
	"github.com/kubideh/kubesearch/search/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

// New returns server App objects.
func New() App {
	return App{
		endpoint:   flag.String("bind-address", ":8080", "IP address and port on which to listen"),
		kubeconfig: kubeconfigFlag(),
	}
}

func kubeconfigFlag() (kubeconfig *string) {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return
}

// App provides everything needed to run kubesearch.
type App struct {
	endpoint   *string
	kubeconfig *string
}

// Run creates a controller that uses the given kubeconfig to watch
// for changes to Kubernetes objects, indexes those objects, and
// provides an API for searching for those objects.
func (a App) Run() {
	flag.Parse()

	// create the Controller to be used by the search API handler
	aController := controller.New(a.client())
	cancel := aController.Start()
	defer cancel()

	mux := http.NewServeMux()
	api.RegisterHandler(mux, index.Searcher(aController.Index()), aController.Store())

	klog.Infoln("Listening on " + *a.endpoint)
	if err := http.ListenAndServe(*a.endpoint, mux); err != nil {
		klog.Fatalln(err)
	}
}

func (a App) client() *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *a.kubeconfig)
	if err != nil {
		klog.Fatalln(err)
	}

	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalln(err)
	}

	return client
}
