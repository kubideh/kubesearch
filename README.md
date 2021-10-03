kubesearch
---

Fulltext search for Kubernetes

## Plan

Build a [custom Kubernetes controller](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md) that indexes the metadata of Kubernetes objects into an [inverted index that supports fulltext search](https://github.com/blevesearch/bleve). The goal is to let the user query for their objects without having to know the exact labels. The result will be a list of object (kind, namespace, name) tuples with which the user can use kubectl to easily get their desired object.

A stretch goal is to build a kubectl plugin that lets formats the API output of kubesearch and enhances the user experience.
