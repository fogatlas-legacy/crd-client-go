/*
Copyright FBK.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/fogatlas/crd-client-go/pkg/apis/fogatlas/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FedFAAppLister helps list FedFAApps.
// All objects returned here must be treated as read-only.
type FedFAAppLister interface {
	// List lists all FedFAApps in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FedFAApp, err error)
	// FedFAApps returns an object that can list and get FedFAApps.
	FedFAApps(namespace string) FedFAAppNamespaceLister
	FedFAAppListerExpansion
}

// fedFAAppLister implements the FedFAAppLister interface.
type fedFAAppLister struct {
	indexer cache.Indexer
}

// NewFedFAAppLister returns a new FedFAAppLister.
func NewFedFAAppLister(indexer cache.Indexer) FedFAAppLister {
	return &fedFAAppLister{indexer: indexer}
}

// List lists all FedFAApps in the indexer.
func (s *fedFAAppLister) List(selector labels.Selector) (ret []*v1alpha1.FedFAApp, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FedFAApp))
	})
	return ret, err
}

// FedFAApps returns an object that can list and get FedFAApps.
func (s *fedFAAppLister) FedFAApps(namespace string) FedFAAppNamespaceLister {
	return fedFAAppNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FedFAAppNamespaceLister helps list and get FedFAApps.
// All objects returned here must be treated as read-only.
type FedFAAppNamespaceLister interface {
	// List lists all FedFAApps in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FedFAApp, err error)
	// Get retrieves the FedFAApp from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FedFAApp, error)
	FedFAAppNamespaceListerExpansion
}

// fedFAAppNamespaceLister implements the FedFAAppNamespaceLister
// interface.
type fedFAAppNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FedFAApps in the indexer for a given namespace.
func (s fedFAAppNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FedFAApp, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FedFAApp))
	})
	return ret, err
}

// Get retrieves the FedFAApp from the indexer for a given namespace and name.
func (s fedFAAppNamespaceLister) Get(name string) (*v1alpha1.FedFAApp, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fedfaapp"), name)
	}
	return obj.(*v1alpha1.FedFAApp), nil
}
