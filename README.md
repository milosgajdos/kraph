**THIS IS a PoC**

[![GoDoc](https://godoc.org/github.com/milosgajdos/kraph?status.svg)](https://godoc.org/github.com/milosgajdos/kraph)
[![Go Report Card](https://goreportcard.com/badge/milosgajdos/kraph)](https://goreportcard.com/report/github.com/milosgajdos/kraph)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Travis CI](https://travis-ci.org/milosgajdos/kraph.svg?branch=master)](https://travis-ci.org/milosgajdos/kraph)

# kraph

`kraph` is an experimental `Go` module that allows to build a graph of API objects. The resulting graph can be represented as [gonum.Graph](https://godoc.org/gonum.org/v1/gonum/graph) which allows for advanced graph analysis!

You can query the resulting `graph`  nodes and edges based on various attributes. Equally, you can also retrieve a subgraph of the graph starting with a chosen node and perform further analysis on it.

At the moment only [kubernetes](https://kubernetes.io/) API object graph is implemented, but the module [hopefully] defines pluggable interfaces which should allow for adding the support for arbitrary APIs, such as AWS etc.

## Getting started

The project provides a simple Makefile which makes basic tasks, such as running tests and building the module simple:

Get dependencies:
```
make dep
```

Run tests:
```shell
make test
```

Build module:
```shell
make build
```

## kctl

There is also a simple command line utility which allows to build and query API object graphs in-memory and display the results.

At the moment it only provides `build` command with `kubernetes/k8s` subcommand which allows to build and query the [kubernetes](https://kubernetes.io/) API object graph.

### HOWTO

Run `kctl` make task to build the `kctl` binary:

Build:
```shell
make kctl
```

`kctl` command line options:
```shell
$ ./kctl help
NAME:
   kctl - build and query API object graphs

USAGE:
   kctl [global options] command [command options] [arguments...]

COMMANDS:
   build    build a graph
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

`kctl` allows to dump the graph in [DOT GraphViz](https://graphviz.gitlab.io/_pages/doc/info/lang.html) format. This can be piped into [GraphViz](https://www.graphviz.org/) tool for further processing. `kctl` is in the pre-alpha state (if it can be called that at all!) and it's more of a debugging help tool at the moment as the focus of the project is the API graph mapping.

**NOTE:** You must have `kubeconfig` properly configured.

```shell
$ ./kctl build k8s -format "dot" | dot -Tsvg > cluster.svg && open cluster.svg
```

**NOTE:** `dot` format is the only available and default format so you can get the same results as above by running the command below, too:
```shell
$ ./kctl build k8s | dot -Tsvg > cluster.svg && open cluster.svg
```
