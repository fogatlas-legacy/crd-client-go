package main

import (
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	clientset "github.com/fogatlas/crd-client-go/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Kubeconfig string `long:"kubeconfig" description:"absolute path to the kubeconfig file."`
	MasterURL  string `long:"master" description:"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster."`
	LogLevel   string `long:"loglevel" description:"The log level." env:"LOGLEVEL" default:"trace"`
}

var opts Options

func main() {
	var parser = flags.NewParser(&opts, flags.Default)
	parser.Parse()
	kubeconfig := opts.Kubeconfig
	masterURL := opts.MasterURL
	configureLog()

	var cfg *rest.Config
	var err error

	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	} else {
		log.Infof("Loading K8S in-cluster configuration")
		cfg, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Fatalf("Unable to get kubeconfig (%s)", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset (%s)", err.Error())
	}

	// get nodes
	nodes, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error: (%s) ", err.Error())
	}
	for _, n := range nodes.Items {
		log.Tracef("Node: (%s)", n.Name)
	}

	crdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building crd clientset (%s)", err.Error())
	}

	// get regions
	regions, err := crdClient.FogatlasV1alpha1().Regions("default").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error: (%s)", err.Error())
	}
	for _, r := range regions.Items {
		log.Tracef("Region: (%s) (%s) (%s) (%s) (%d)", r.Spec.Id, r.Spec.Name, r.Spec.Description,
			r.Spec.Location, r.Spec.Tier)
	}
}

func configureLog() {
	var loggerLevel log.Level
	logLevel := opts.LogLevel
	switch logLevel {
	case "trace":
		loggerLevel = log.TraceLevel
	case "debug":
		loggerLevel = log.DebugLevel
	case "info":
		loggerLevel = log.InfoLevel
	case "warn":
		loggerLevel = log.WarnLevel
	case "error":
		loggerLevel = log.ErrorLevel
	case "fatal":
		loggerLevel = log.FatalLevel
	case "panic":
		loggerLevel = log.PanicLevel
	default:
		loggerLevel = log.InfoLevel
	}

	log.SetLevel(loggerLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

}
