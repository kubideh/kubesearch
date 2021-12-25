package app

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

// newFlags returns the immutableFlags for App.
func newFlags() immutableFlags {
	return immutableFlags{
		bindAddress: flag.String("bind-address", ":8080", "IP address and port on which to listen"),
		kubeConfig:  kubeConfigFlag(),
	}
}

func kubeConfigFlag() (kubeConfig *string) {
	// It's convention to use `kubeconfig` as the flag name.
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return
}

// immutableFlags is a collection of flags used to configure the
// App. Each flag will be populated with values from the command-
// line after calling Parse().
type immutableFlags struct {
	bindAddress *string // bindAddress is an address that can be used by `http.ListenAndServe`
	kubeConfig  *string // kubeConfig is a path string that can be used to create Kubernetes clients
}

// BindAddress returns an address that can be used by
// `http.ListenAndServe`, and it's populated by a value from the
// command-line.
func (f immutableFlags) BindAddress() string {
	return *f.bindAddress
}

// KubeConfig returns a path string that can be used to create
// Kubernetes clients, and it's populated by a value from the
// command-line.
func (f immutableFlags) KubeConfig() string {
	return *f.kubeConfig
}

// Parse populates this collection of immutableFlags with values from the
// command-line.
func (f immutableFlags) Parse() {
	flag.Parse()
}
