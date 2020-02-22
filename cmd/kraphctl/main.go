package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/milosgajdos83/kraph"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	var (
		kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster")
		masterURL  = flag.String("master", "", "The URL of the Kubernetes API server")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// set up signal handler
	sigChan := setupSignalHandler()

	config, err := GetKubeConfig(*masterURL, *kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes clientset: %w", err)
	}

	k, err := kraph.New(client)
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	if err := k.Build(); err != nil {
		return fmt.Errorf("failed to build kraph: %w", err)
	}

	select {
	case s := <-sigChan:
		fmt.Fprintf(stdout, "caught %s signal, terminating\n", s)
		return nil
	default:
		// do nothing for now
	}

	return nil
}

// GetKubeConfig builds kubernetes configuration and returns it.
// It looks for kubernetes config file in the following order:
// 	1. kubeconfig
// 	2. $KUBECONFIG environment variable
// 	3. $HOMEDIR/.kube/config
// It returns error if the configuration could not be built.
func GetKubeConfig(masterURL, kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}
		}
	}

	// NOTE: if neither masterURL nor kubeconfig is provided this defaults to in-cluster config
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed building kubernetes config: %v", err)
	}

	return config, nil
}

// setupSignalHandler makes signal handler for catching os.Interrupt and returns it
// NOTE: we could potentially expand this to variadic number of signals, but meh for now
func setupSignalHandler() <-chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	return signalChan
}
