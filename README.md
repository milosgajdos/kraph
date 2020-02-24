**THIS IS a PoC**

# kraph

`kraph` is a `Go` package which allows to build a graph of all available Kubernetes API resources and encodes them into [gonum.graph](https://godoc.org/gonum.org/v1/gonum/graph)

# HOWTO

**NOTE:** You must have `kubeconfig` properly configured

Get all dependencies:
```shell
go get
```

Run tests
```shell
go test
```

Build kraphctl:
```shell
go build cmd/kraphctl/main.go -o kraphctl
```

Run `kraphctl`:
```shell
./kraphctl | dot -Tsvg > cluster.svg && open cluster.svg
```

# Example

vanilla [kind](https://kind.sigs.k8s.io/) cluster:

![Kind cluster](examples/kind.svg)
