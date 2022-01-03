package app

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

// CreateImmutableServerFlags returns the ImmutableServerFlags for
// App. A list of the flags and their defaults are now given.
//
// -bind-address (default: :8080)
// -kubeconfig (default $HOME/.kube/config if $HOME is set; empty string otherwise)
func CreateImmutableServerFlags() ImmutableServerFlags {
	return CreateImmutableServerFlagsWithBindAddress(":8080")
}

// CreateImmutableServerFlagsWithBindAddress returns the
// ImmutableServerFlags for App, and it uses the given bindAddress
// as a default value for the flag `-bind-address`.
func CreateImmutableServerFlagsWithBindAddress(bindAddress string) ImmutableServerFlags {
	return ImmutableServerFlags{
		bindAddress: flag.String("bind-address", bindAddress, "IP address and port on which to listen"),
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

// ImmutableServerFlags is a collection of flags used to configure
// the App. Each flag will be populated with values from the
// command-line after calling Parse().
type ImmutableServerFlags struct {
	bindAddress *string // bindAddress is an address that can be used by `http.ListenAndServe`
	kubeConfig  *string // kubeConfig is a path string that can be used to create Kubernetes clients
}

// BindAddress returns an address that can be used by
// `http.ListenAndServe`, and it's populated by a value from the
// command-line.
func (f ImmutableServerFlags) BindAddress() string {
	return *f.bindAddress
}

// KubeConfig returns a path string that can be used to create
// Kubernetes clients, and it's populated by a value from the
// command-line.
func (f ImmutableServerFlags) KubeConfig() string {
	return *f.kubeConfig
}

// Parse populates this collection of ImmutableServerFlags with values from the
// command-line.
func (f ImmutableServerFlags) Parse() {
	flag.Parse()
}
