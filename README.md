[![CI Workflow](https://github.com/kubideh/kubesearch/actions/workflows/main.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/main.yml)
[![CodeQL Analysis](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kubideh/kubesearch.svg)](https://pkg.go.dev/github.com/kubideh/kubesearch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubideh/kubesearch)](https://goreportcard.com/report/github.com/kubideh/kubesearch)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![codecov](https://codecov.io/gh/kubideh/kubesearch/branch/main/graph/badge.svg?token=YP1EDH6PTH)](https://codecov.io/gh/kubideh/kubesearch)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# kubesearch

<img src="https://github.com/kubernetes/community/blob/master/icons/png/control_plane_components/labeled/api-256.png?raw=true" width="100">

----

Fulltext search for Kubernetes API objects

The Kubernetes API supports exact matches in order to lookup API
objects. For example, the user can fetch a Pod by name or list Pods
using labels. In both scenarios, the user must know exactly the name
of the Pod or labels of the Pod.

Kubesearch and the client kubectl-search let the user **search**
for relevant API objects without having to know the exact name,
namespace, or labels.

## Getting started

### Containerized server
Include steps to pull a container or install a release.

```console
```

### Build client and server from source

```console
go install ./...
```

### Pull client and server from homebrew

## Usage

### Run kubesearch as a stand-alone service

```console
kubesearch
```

### Search for Kubernetes objects using kubectl

```console
kubectl create ns flargle
kubectl run blargle -n flargle --image=nginx:alpine
kubectl run boggle -n flargle --image=nginx:latest
kubectl search blargle
kubectl search flargle
kubectl search nginx
kubectl search \"nginx:alpine\"
```

## API

`/v1/search?query=<fulltext query string>` # Search using a phrase query by surrounding the query in `"` (quotes)

## To do for v1.0.0

1. Develop a better tokenizer
  - [ANTLR4](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
  - [Unicode Text Segmentation] (https://unicode.org/reports/tr29/)
2. Normalize terms to lowercase
3. Support phrase-search (searching for exact phrases which may include token separators)
4. Index annotations, container names, images names, labels, and volume names
5. Index DaemonSets
6. Index arbitrary fields
7. Make indexable fields configurable
8. Index arbitrary resource types
9. Use a treap
10. Consider supporting configurable policies in order to control access to the API (OPA)
11. Consider vector space model for retrieval
12. Release using homebrew
13. Metrics

## To do for v2.0

1. Support kubesearch as an API extension
2. Add a client that searches using the API extension
3.
