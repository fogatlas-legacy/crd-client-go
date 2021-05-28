package v1alpha1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/fogatlas/crd-client-go/pkg/apis/fogatlas/v1alpha1"
	fadeplclientfake "github.com/fogatlas/crd-client-go/pkg/generated/clientset/versioned/fake"
	fadeplinformers "github.com/fogatlas/crd-client-go/pkg/generated/informers/externalversions"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

type fixture struct {
	fadeplClientSet       *fadeplclientfake.Clientset
	kubeClient            *k8sfake.Clientset
	fadeplInformerFactory fadeplinformers.SharedInformerFactory
	kubeInformerFactory   kubeinformers.SharedInformerFactory
}

func newFixture() *fixture {
	f := &fixture{}
	objects := []runtime.Object{}
	f.fadeplClientSet = fadeplclientfake.NewSimpleClientset(objects...)
	f.kubeClient = k8sfake.NewSimpleClientset(objects...)
	f.fadeplInformerFactory = fadeplinformers.NewSharedInformerFactory(f.fadeplClientSet, time.Second*0)
	f.kubeInformerFactory = kubeinformers.NewSharedInformerFactory(f.kubeClient, time.Second*0)
	return f
}

func int32Ptr(i int32) *int32 { return &i }

func newDeployment(name string, replicas *int32, cpuRequested string, memoryRequested string) apps.Deployment {
	selector := make(map[string]string)
	selector["name"] = name
	resourceList := make(map[v1.ResourceName]resource.Quantity)
	resourceList[v1.ResourceCPU] = resource.MustParse(cpuRequested)
	resourceList[v1.ResourceMemory] = resource.MustParse(memoryRequested)
	d := apps.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   metav1.NamespaceDefault,
			Annotations: make(map[string]string),
		},
		Spec: apps.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selector,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  name,
							Image: "nginx",
							Resources: v1.ResourceRequirements{
								Requests: resourceList,
							},
						},
					},
				},
			},
		},
	}
	return d
}

func completeFADepl(name string, replicas *int32) *v1alpha1.FADepl {
	return &v1alpha1.FADepl{
		TypeMeta: metav1.TypeMeta{APIVersion: v1alpha1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Finalizers: []string{"finalizer.fogatlas.fbk.eu"},
			Namespace:  metav1.NamespaceDefault,
		},
		Spec: v1alpha1.FADeplSpec{
			ExternalEndpoints: []string{"cam1", "cam2", "broker1", "sens1"},
			Algorithm:         "DAG",
			Microservices: []*v1alpha1.FADeplMicroservice{
				{
					Name:       "processor",
					Deployment: newDeployment("processor", replicas, "100m", "400M"),
				},
				{
					Name: "driver1",
					Regions: []*v1alpha1.FARegion{
						{
							RegionRequired: "003-003",
						},
					},
					Deployment: newDeployment("driver1", replicas, "100m", "100M"),
				},
				{
					Name: "driver2",
					Regions: []*v1alpha1.FARegion{
						{
							RegionRequired: "003-003",
						},
					},
					Deployment: newDeployment("driver2", replicas, "100m", "100M"),
				},
				{
					Name:       "driver3",
					Deployment: newDeployment("driver3", replicas, "100m", "100M"),
				},
				{
					Name:       "analytics",
					Deployment: newDeployment("analytics", replicas, "100m", "400M"),
				},
			},
			DataFlows: []*v1alpha1.FADeplDataFlow{
				{
					BandwidthRequired: resource.MustParse("5M"),
					LatencyRequired:   resource.MustParse("20"),
					SourceID:          "cam1",
					DestinationID:     "driver1",
				},
				{
					BandwidthRequired: resource.MustParse("5M"),
					LatencyRequired:   resource.MustParse("20"),
					SourceID:          "cam2",
					DestinationID:     "driver2",
				},
				{
					BandwidthRequired: resource.MustParse("5M"),
					LatencyRequired:   resource.MustParse("20"),
					SourceID:          "broker1",
					DestinationID:     "driver3",
				},
				{
					BandwidthRequired: resource.MustParse("100k"),
					LatencyRequired:   resource.MustParse("500"),
					SourceID:          "driver1",
					DestinationID:     "processor",
				},
				{
					BandwidthRequired: resource.MustParse("100k"),
					LatencyRequired:   resource.MustParse("500"),
					SourceID:          "driver2",
					DestinationID:     "processor",
				},
				{
					BandwidthRequired: resource.MustParse("100k"),
					LatencyRequired:   resource.MustParse("500"),
					SourceID:          "driver3",
					DestinationID:     "processor",
				},
				{
					BandwidthRequired: resource.MustParse("100k"),
					LatencyRequired:   resource.MustParse("500"),
					SourceID:          "processor",
					DestinationID:     "analytics",
				},
				{
					BandwidthRequired: resource.MustParse("100k"),
					LatencyRequired:   resource.MustParse("500"),
					SourceID:          "analytics",
					DestinationID:     "sens1",
				},
			},
		},
	}
}

func emptyFADepl(name string, replicas *int32) *v1alpha1.FADepl {
	return &v1alpha1.FADepl{
		TypeMeta: metav1.TypeMeta{APIVersion: v1alpha1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
	}
}

// TestSubmitFADepl creates a complete and an empty fadepl and verify that the
// created ones are identical to the ones submitted to Create().
func TestSubmitFADepl(t *testing.T) {
	f := newFixture()

	var tests = []struct {
		testName string
		fadepl   *v1alpha1.FADepl
	}{
		{
			testName: "CompleteFADepl",
			fadepl:   completeFADepl("CompleteFADepl", int32Ptr(1)),
		},
		{
			testName: "EmptyFADepl",
			fadepl:   emptyFADepl("EmptyFADepl", int32Ptr(1)),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			fadepl, err := f.fadeplClientSet.FogatlasV1alpha1().FADepls(test.fadepl.Namespace).Create(context.TODO(), test.fadepl, metav1.CreateOptions{})
			if err != nil {
				t.Error(err)
			} else {
				valid := cmp.Equal(fadepl, test.fadepl, nil)
				if !valid {
					diff := cmp.Diff(fadepl, test.fadepl, nil)
					t.Errorf("Expected and actual differ: %s", diff)
				}
			}
		})
	}
}

// TestGetFADepl submits a fadepl through the Informer and gets a complete fadepl in two ways:
// - through the Lister
// - through the clientSet
// In both cases verifies of the retrieved one is identical to the template.
func TestGetFADepl(t *testing.T) {
	f := newFixture()

	fadepl := completeFADepl("CompleteFADepl", int32Ptr(1))
	// add directly to the store without passing through a Create call
	f.fadeplInformerFactory.Fogatlas().V1alpha1().FADepls().Informer().GetIndexer().Add(fadepl)

	// reading from the lister
	lister := f.fadeplInformerFactory.Fogatlas().V1alpha1().FADepls().Lister()
	gotFADepl, err := lister.FADepls(fadepl.Namespace).Get(fadepl.Name)
	if err != nil {
		t.Error(err)
	} else {
		valid := cmp.Equal(gotFADepl, fadepl, nil)
		if !valid {
			diff := cmp.Diff(gotFADepl, fadepl, nil)
			t.Errorf("Expected and actual differ: %s", diff)
		}
	}
	// writing and reading from clienset
	_, err = f.fadeplClientSet.FogatlasV1alpha1().FADepls(metav1.NamespaceDefault).Create(context.TODO(), fadepl, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
	gotFADepl, err = f.fadeplClientSet.FogatlasV1alpha1().FADepls(fadepl.Namespace).Get(context.TODO(), fadepl.Name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	} else {
		valid := cmp.Equal(gotFADepl, fadepl, nil)
		if !valid {
			diff := cmp.Diff(gotFADepl, fadepl, nil)
			t.Errorf("Expected and actual differ: %s", diff)
		}
	}
}

// TestUpdateFADepl updates a FADepl
func TestUpdateFADepl(t *testing.T) {
	f := newFixture()
	fadepl := completeFADepl("CompleteFADepl", int32Ptr(1))
	_, err := f.fadeplClientSet.FogatlasV1alpha1().FADepls(fadepl.Namespace).Create(context.TODO(), fadepl, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
	} else {
		gotFADepl, err := f.fadeplClientSet.FogatlasV1alpha1().FADepls(fadepl.Namespace).Get(context.TODO(), fadepl.Name, metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}
		faDeplCopy := gotFADepl.DeepCopy()
		_, err = f.fadeplClientSet.FogatlasV1alpha1().FADepls(faDeplCopy.Namespace).Update(context.TODO(), faDeplCopy, metav1.UpdateOptions{})
		if err != nil {
			t.Error(err)
		}
	}
}
