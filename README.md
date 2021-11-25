[![CI Workflow](https://github.com/kubideh/kubesearch/actions/workflows/main.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/main.yml)
[![CodeQL Analysis](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kubideh/kubesearch/actions/workflows/codeql-analysis.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kubideh/kubesearch.svg)](https://pkg.go.dev/github.com/kubideh/kubesearch)

# kubesearch

<img src="https://github.com/kubernetes/community/blob/master/icons/png/control_plane_components/labeled/api-256.png?raw=true" width="100">

----

Fulltext search for Kubernetes

## Requirements

```console
brew install go
brew install goreleaser
```

## Getting started

```console
go test -v ./...
```

## Plan

Build a custom Kubernetes object reflector which indexes the
metadata of Kubernetes objects into an inverted index that supports
fulltext search. The goal is to let the user query for their
objects without having to know the exact labels. The result will be
a list of object (kind, namespace, name) tuples with which the user
can use kubectl to easily get their desired object.

A stretch goal is to build a kubectl plugin that lets formats the
API output of kubesearch and enhances the user experience.

## API

`/v1/search?query=`
