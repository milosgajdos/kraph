package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/milosgajdos/kraph"
	"github.com/milosgajdos/kraph/api/k8s"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/memory"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	api        string
	kubeconfig string
	master     string
	namespace  string
	format     string
	graphStore string
)

// K8s returns K8s subcommand for build command
func K8s() *cli.Command {
	return &cli.Command{
		Name:     "kubernetes",
		Aliases:  []string{"k8s"},
		Category: "build",
		Usage:    "kubernetes graph",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api",
				Aliases:     []string{"a"},
				Value:       "all",
				Usage:       "kubernetes APIs",
				Destination: &api,
			},
			&cli.StringFlag{
				Name:        "store",
				Aliases:     []string{"s"},
				Value:       "memory",
				Usage:       "graph store",
				Destination: &graphStore,
			},
			&cli.StringFlag{
				Name:        "kubeconfig",
				Aliases:     []string{"c"},
				Usage:       "Path to a kubeconfig",
				Destination: &kubeconfig,
			},
			&cli.StringFlag{
				Name:        "master",
				Aliases:     []string{"m"},
				Usage:       "URL of the Kubernetes API server",
				Destination: &master,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"ns"},
				Usage:       "Kubernetes namespace",
				Destination: &namespace,
			},
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Value:       "dot",
				Usage:       "print graph in a given format",
				Destination: &format,
			},
		},
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}
}

// getKubeConfig builds kubernetes configuration and returns it.
// It looks for kubernetes config file in the following order:
// 	1. kubeconfig
// 	2. $KUBECONFIG environment variable
// 	3. $HOMEDIR/.kube/config
// It returns error if the configuration could not be built.
func getKubeConfig(masterURL, kubeconfig string) (*rest.Config, error) {
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

func graphToOut(g store.Graph, format string) (string, error) {
	switch format {
	case "dot":
		dotGraph := g.(store.DOTGraph)
		return dotGraph.DOT()
	default:
		dotGraph := g.(store.DOTGraph)
		return dotGraph.DOT()
	}
}

func run(ctx *cli.Context) error {
	config, err := getKubeConfig(master, kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	// adjust configuration for faster scan
	config.QPS = 100
	config.Burst = 100

	discClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes clientset: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes dynamic client: %w", err)
	}

	var gstore store.Store
	storeID := "kctl"

	switch graphStore {
	case "memory":
		gstore, err = memory.NewStore(storeID)
		if err != nil {
			return err
		}
	default:
		gstore, err = memory.NewStore(storeID)
		if err != nil {
			return err
		}
	}

	k, err := kraph.New(kraph.Store(gstore))
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	_, err = k.Build(k8s.NewClient(discClient.Discovery(), dynClient, ctx.Context, k8s.Namespace(namespace)))
	if err != nil {
		return fmt.Errorf("failed to build kraph: %w", err)
	}

	graphOut, err := graphToOut(k.Store(), format)
	if err != nil {
		return err
	}

	fmt.Println(graphOut)

	return nil
}
