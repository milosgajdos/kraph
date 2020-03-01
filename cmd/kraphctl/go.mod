module github.com/milosgajdos/kraph/cmd/kraphctl

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/milosgajdos/kraph v0.0.0-20200229111353-cd6d608c0b31
	k8s.io/client-go v0.17.3
)

replace github.com/milosgajdos/kraph => /Users/milosgajdos/go/src/github.com/milosgajdos/kraph
