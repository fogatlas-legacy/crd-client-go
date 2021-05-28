package main

import (
	"context"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	clientset "github.com/fogatlas/crd-client-go/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Options collects the possibile command line options
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
	log.Tracef("Getting Nodes")
	nodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
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
	log.Tracef("Getting Regions")
	regions, err := crdClient.FogatlasV1alpha1().Regions("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error: (%s)", err.Error())
	}
	for _, r := range regions.Items {
		log.Tracef("Region: (%s) (%s) (%s) (%s) (%d)", r.Spec.ID, r.Spec.Name, r.Spec.Description,
			r.Spec.Location, r.Spec.Tier)
	}

	// get fadepl
	log.Tracef("Getting FADepls")
	fadepls, err := crdClient.FogatlasV1alpha1().FADepls("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error: (%s)", err.Error())
	}
	for _, fa := range fadepls.Items {
		log.Tracef("JSON Fadepl : (%#v)", fa)
		for _, ms := range fa.Spec.Microservices {
			log.Tracef("MS name is (%s); replicas is (%d)", ms.Deployment.Name, *ms.Deployment.Spec.Replicas)
		}
	}

	//get list of crds
	log.Tracef("Getting CRDs")
	metaClient, err := metadata.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building metadata clientset (%s)", err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	crds, err := metaClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error: (%s) ", err.Error())
	}
	for _, crd := range crds.Items {
		log.Tracef("CRD: (%s)", crd.Name)
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
