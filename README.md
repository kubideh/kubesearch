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

## Requirements

```console
brew install go
```

## Getting started

```console
go install ./...
```

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

## To do

1. Don't store duplicate postings in the same postings list; increment term frequency instead
2. Parse queries into terms and combine results
3. Basic ranked retrieval by term frequency
4. Combine results to support exact phrases
5. Index annotations, container names, images names, labels, and volume names

