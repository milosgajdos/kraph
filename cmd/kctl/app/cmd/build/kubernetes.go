package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/milosgajdos/kraph"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/k8s/owner"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/memory"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig string
	master     string
	namespace  string
	format     string
	graphStore string
	storeURL   string
	storeID    string
	graphType  string
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
				Name:        "store",
				Aliases:     []string{"s"},
				Value:       "memory",
				Usage:       "graph store",
				Destination: &graphStore,
			},
			&cli.StringFlag{
				Name:        "store-id",
				Aliases:     []string{"id"},
				Value:       "kctl",
				Usage:       "store ID",
				Destination: &storeID,
			},
			&cli.StringFlag{
				Name:        "store-url",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "URL of the store",
				EnvVars:     []string{"STORE_URL"},
				Destination: &storeURL,
			},
			&cli.StringFlag{
				Name:        "graph",
				Aliases:     []string{"g"},
				Value:       "owner",
				Usage:       "type of graph",
				Destination: &graphType,
			},
			&cli.StringFlag{
				Name:        "kubeconfig",
				Aliases:     []string{"c"},
				Usage:       "path to kubeconfig",
				Destination: &kubeconfig,
			},
			&cli.StringFlag{
				Name:        "master",
				Aliases:     []string{"m"},
				Usage:       "URL of kubernetes API server",
				Destination: &master,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"ns"},
				Usage:       "kubernetes namespace",
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
		return "", fmt.Errorf("unsupported output format: %s", format)
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

	switch graphStore {
	case "memory":
		gstore, err = memory.NewStore(storeID, store.Options{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported store: %s", graphStore)
	}

	k, err := kraph.New(gstore)
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	var filters []kraph.Filter

	var client api.Client

	switch graphType {
	case "owner":
		client = owner.NewClient(ctx.Context, discClient.Discovery(), dynClient, owner.Namespace(namespace))
	default:
		return fmt.Errorf("unsupported graph type: %s", graphType)
	}

	if err = k.Build(client, filters...); err != nil {
		return fmt.Errorf("failed to build kraph: %w", err)
	}

	graphOut, err := graphToOut(k.Store().Graph(), format)
	if err != nil {
		return err
	}

	fmt.Println(graphOut)

	return nil
}
