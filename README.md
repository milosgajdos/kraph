**THIS IS a PoC**

[![GoDoc](https://godoc.org/github.com/milosgajdos/kraph?status.svg)](https://godoc.org/github.com/milosgajdos/kraph)
[![Go Report Card](https://goreportcard.com/badge/milosgajdos/kraph)](https://goreportcard.com/report/github.com/milosgajdos/kraph)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Travis CI](https://travis-ci.org/milosgajdos/kraph.svg?branch=master)](https://travis-ci.org/milosgajdos/kraph)

# kraph

`kraph` is an experimental `Go` module which allows to build a graph of API objects. The resulting graph can be represented as [gonum.Graph](https://godoc.org/gonum.org/v1/gonum/graph) which allows for advanced graph analysis!

You can query the resulting `graph`  nodes and edges based on various attributes. Equally, you can also retrieve a subgraph of a chosen node and perform further analysis on it.

At the moment only [kubernetes](https://kubernetes.io/) API object graph is implemented, but the module defines pluggable interfaces which should allow for expanding the support for arbitrary API objects, such as AWS etc.

## kraphctl

The project provides an example cli utility which demonstrates how the kraph can be built. In this particular case it demonstrates it on using the kubernetes API.

`kraphctl` is lets you build the graph of the the Kubernetes API objects and then dump the resulting graph in [DOT GraphViz](https://graphviz.gitlab.io/_pages/doc/info/lang.html) format. This can be piped into the [GraphViz](https://www.graphviz.org/) tool for further processing.interact with the `kraph`. It's in the pre-alpha state (if it can be called that at all!).

### HOWTO

**NOTE:** You must have `kubeconfig` properly configured

There is a simple Makefile which makes basic tasks, such as running tests and building the module simple:

Build:
```shell
make kraphctl
```

`kraphctl` command line options:
```shell
$ ./kraphctl -h
Usage of ./kraphctl:
  -kubeconfig string
    	Path to a kubeconfig. Only required if out-of-cluster
  -master string
    	The URL of the Kubernetes API server
  -namespace string
    	Kubernetes namespace
```

Run `kraphctl`:
```shell
./kraphctl | dot -Tsvg > cluster.svg && open cluster.svg
```
