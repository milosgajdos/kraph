package kraph

import (
	"testing"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewKraph(t *testing.T) {
	_, err := New(testclient.NewSimpleClientset())
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}
}

func TestBuild(t *testing.T) {
	k, err := New(testclient.NewSimpleClientset())
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	if err := k.Build(); err != nil {
		t.Errorf("failed to build kraph: %v", err)
	}
}
