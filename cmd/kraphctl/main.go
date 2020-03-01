package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/milosgajdos/kraph"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	signalChan := setupSignalHandler()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(1)
	}()

	if err := run(ctx, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// parse cli flags
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	var (
		kubeconfig = flags.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster")
		masterURL  = flags.String("master", "", "The URL of the Kubernetes API server")
		namespace  = flags.String("namespace", "", "Kubernetes namespace")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	config, err := GetKubeConfig(*masterURL, *kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	// adjust configuration for faster scan
	config.QPS = 100
	config.Burst = 100

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes clientset: %w", err)
	}

	clientDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes dynamic client: %w", err)
	}

	k, err := kraph.New(client.Discovery(), clientDynamic, kraph.Namespace(*namespace))
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	if err := k.Build(ctx); err != nil {
		return fmt.Errorf("failed to build kraph: %w", err)
	}

	dotKraph, err := k.DOT()
	//_, err = k.DOT()
	if err != nil {
		log.Fatal(err)
		//return err
	}

	fmt.Println(dotKraph)

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
func setupSignalHandler() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	return signalChan
}
