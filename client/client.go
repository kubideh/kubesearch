// Package client provides access to a global K8s client.
package client

import "k8s.io/client-go/kubernetes"

var singleton kubernetes.Interface

// Client returns the K8s client.
func Client() kubernetes.Interface {
	return singleton
}

// SetClient replaces the K8s client.
func SetClient(client kubernetes.Interface) {
	singleton = client
}
