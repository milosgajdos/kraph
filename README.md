**THIS IS a PoC**

[![GoDoc](https://godoc.org/github.com/milosgajdos/kraph?status.svg)](https://godoc.org/github.com/milosgajdos/kraph)
[![Go Report Card](https://goreportcard.com/badge/milosgajdos/kraph)](https://goreportcard.com/report/github.com/milosgajdos/kraph)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Travis CI](https://travis-ci.org/milosgajdos/kraph.svg?branch=master)](https://travis-ci.org/milosgajdos/kraph)

# kraph

`kraph` is an experimental `Go` module which allows to build a graph of all Kubernetes API objects. The resulting graph is represented as [gonum.Graph](https://godoc.org/gonum.org/v1/gonum/graph) which allows for advanced graph analysis!

You can query the resulting `graph`  nodes and edges based on various attributes. Equally, you can also retrieve a subgraph of a chosen node and perform further analysis on it.

# HOWTO

**NOTE:** You must have `kubeconfig` properly configured

There is a simple Makefile which makes basic tasks, such as running tests and building the module simple:

Run tests:
```shell
make test
```

Build module:
```shell
make build
```

Build `kraphctl`:
```shell
make kraphctl
```

## kraphctl

`kraphctl` is a simple command line utility which lets you query the Kubernetes API and dump the graph of all discovered API objects in [DOT GraphViz](https://graphviz.gitlab.io/_pages/doc/info/lang.html) format. This can be piped into the [GraphViz](https://www.graphviz.org/) tool for further processing.interact with the `kraph`. It's in the pre-alpha state (if it can be called that at all!).

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
