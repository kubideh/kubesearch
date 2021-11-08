![CI](https://github.com/kubideh/kubesearch/actions/workflows/main.yml/badge.svg)

# kubesearch

<img src="https://github.com/kubernetes/community/blob/master/icons/png/control_plane_components/labeled/api-256.png?raw=true" width="100">

----

Fulltext search for Kubernetes

## Requirements

```console
brew install go
brew install goreleaser
go get github.com/onsi/ginkgo/ginkgo
```

## Getting started

```console
ginkgo -v ./...
```

## Plan

Build a [custom Kubernetes controller](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md) that indexes the metadata of Kubernetes objects into an [inverted index that supports fulltext search](https://github.com/blevesearch/bleve). The goal is to let the user query for their objects without having to know the exact labels. The result will be a list of object (kind, namespace, name) tuples with which the user can use kubectl to easily get their desired object.

A stretch goal is to build a kubectl plugin that lets formats the API output of kubesearch and enhances the user experience.

## API

`/v1/search?query=`
