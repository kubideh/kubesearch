kubesearch
---

[![CI Workflow](https://github.com/kubideh/kubesearch/actions/workflows/main.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/main.yml)
[![CodeQL Analysis](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kubideh/kubesearch.svg)](https://pkg.go.dev/github.com/kubideh/kubesearch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubideh/kubesearch)](https://goreportcard.com/report/github.com/kubideh/kubesearch)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![codecov](https://codecov.io/gh/kubideh/kubesearch/branch/main/graph/badge.svg?token=YP1EDH6PTH)](https://codecov.io/gh/kubideh/kubesearch)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Kubesearch** is fulltext search for Kubernetes.

## Introduction

<img src="https://github.com/kubernetes/community/blob/master/icons/png/control_plane_components/labeled/api-256.png?raw=true" width="100">

The Kubernetes API supports exact matches in order to lookup API
objects. For example, the user can fetch a Pod by name or list Pods
using labels. In both scenarios, the user must know exactly the name
of the Pod or labels of the Pod. In all scenarios, the user must
know exactly the name or labels of any object being searched for.

Kubesearch and the client kubectl-search let the user **search**
for relevant API objects without having to know the exact name,
namespace, or labels. The results are formatted as commands that
can be copied and executed in order to retrieve the desired API
object.

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

1. Release using homebrew
2. Metrics
3. Develop a better tokenizer
4. Normalize terms to lowercase
5. Support phrase-search (searching for exact phrases which may include token separators)
6. Index annotations, container names, images names, labels, and volume names
7. Index DaemonSets
8. Make indexable resource types and fields configurable
9. Index arbitrary fields
10. Index arbitrary resource types
11. Use a treap
12. Consider vector space model for retrieval

## To do for v2.0

1. Support kubesearch as an API extension
2. Add a client that searches using the API extension
3. Consider supporting configurable policies in order to control access to the API (OPA)
4. Format results of kubectl-search in order to copy and paste the results as an executable command

## References

* [Introduction to Information Retrieval](https://nlp.stanford.edu/IR-book/information-retrieval-book.html)
* [Unicode Text Segmentation](https://unicode.org/reports/tr29/)
* [ANTLR4](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
* [Faster and smaller inverted indices with treaps](https://dl.acm.org/doi/10.1145/2484028.2484088)

