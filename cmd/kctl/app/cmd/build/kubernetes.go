package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"

	"github.com/milosgajdos/kraph"
	"github.com/milosgajdos/kraph/pkg/api/k8s"
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
	kinds      string
	kubeconfig string
	master     string
	namespace  string
	format     string
	graphStore string
	storeURL   string
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
				Name:        "kinds",
				Aliases:     []string{"k"},
				Value:       "all",
				Usage:       "filter by resource kinds (comma separated)",
				Destination: &kinds,
			},
			&cli.StringFlag{
				Name:        "store",
				Aliases:     []string{"s"},
				Value:       "memory",
				Usage:       "graph store",
				Destination: &graphStore,
			},
			&cli.StringFlag{
				Name:        "store-url",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "URL of a remote graph store",
				EnvVars:     []string{"STORE_URL"},
				Destination: &storeURL,
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
		gstore, err = memory.NewStore(storeID, store.Options{})
		if err != nil {
			return err
		}
	default:
		gstore, err = memory.NewStore(storeID, store.Options{})
		if err != nil {
			return err
		}
	}

	k, err := kraph.New(kraph.Store(gstore))
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	var filters []kraph.Filter
	if len(kinds) > 0 && kinds != "all" {
		for _, kind := range strings.Split(kinds, ",") {
			filters = append(filters,
				func(object api.Object) bool { return object.Resource().Kind() == kind },
			)
		}
	}

	client := k8s.NewClient(ctx.Context, discClient.Discovery(), dynClient, k8s.Namespace(namespace))

	// TODO: Build now returns store.Graph
	// there is no need to call k.Store() as below
	_, err = k.Build(client, filters...)
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
