**THIS IS a PoC**

[![GoDoc](https://godoc.org/github.com/milosgajdos/kraph?status.svg)](https://godoc.org/github.com/milosgajdos/kraph)
[![Go Report Card](https://goreportcard.com/badge/milosgajdos/kraph)](https://goreportcard.com/report/github.com/milosgajdos/kraph)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Travis CI](https://travis-ci.org/milosgajdos/kraph.svg?branch=master)](https://travis-ci.org/milosgajdos/kraph)

# kraph

`kraph` is an experimental `Go` module which allows to build a graph of all Kubernetes API resources as supported by its API server. The graph is encoded into [gonum.Graph](https://godoc.org/gonum.org/v1/gonum/graph) which allows for advanced graph analysis!

The project also provides a simple tool which queries the Kubernetes API and dumps the graph of all discovered API resources in [DOT GraphViz](https://graphviz.gitlab.io/_pages/doc/info/lang.html) format. This can be piped into the [GraphViz](https://www.graphviz.org/) tool for further processing.

# HOWTO

**NOTE:** You must have `kubeconfig` properly configured

Get all dependencies:
```shell
go get
```

Run tests:
```shell
go test
```

Build `kraphctl`:
```shell
go build cmd/kraphctl/main.go -o kraphctl
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

## Example

vanilla [kind](https://kind.sigs.k8s.io/) cluster:

![Kind cluster](examples/kind.svg)
