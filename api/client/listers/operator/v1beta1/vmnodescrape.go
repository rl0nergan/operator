/*


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
// Code generated by lister-gen-v0.30. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VMNodeScrapeLister helps list VMNodeScrapes.
// All objects returned here must be treated as read-only.
type VMNodeScrapeLister interface {
	// List lists all VMNodeScrapes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.VMNodeScrape, err error)
	// VMNodeScrapes returns an object that can list and get VMNodeScrapes.
	VMNodeScrapes(namespace string) VMNodeScrapeNamespaceLister
	VMNodeScrapeListerExpansion
}

// vMNodeScrapeLister implements the VMNodeScrapeLister interface.
type vMNodeScrapeLister struct {
	indexer cache.Indexer
}

// NewVMNodeScrapeLister returns a new VMNodeScrapeLister.
func NewVMNodeScrapeLister(indexer cache.Indexer) VMNodeScrapeLister {
	return &vMNodeScrapeLister{indexer: indexer}
}

// List lists all VMNodeScrapes in the indexer.
func (s *vMNodeScrapeLister) List(selector labels.Selector) (ret []*v1beta1.VMNodeScrape, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VMNodeScrape))
	})
	return ret, err
}

// VMNodeScrapes returns an object that can list and get VMNodeScrapes.
func (s *vMNodeScrapeLister) VMNodeScrapes(namespace string) VMNodeScrapeNamespaceLister {
	return vMNodeScrapeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VMNodeScrapeNamespaceLister helps list and get VMNodeScrapes.
// All objects returned here must be treated as read-only.
type VMNodeScrapeNamespaceLister interface {
	// List lists all VMNodeScrapes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.VMNodeScrape, err error)
	// Get retrieves the VMNodeScrape from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.VMNodeScrape, error)
	VMNodeScrapeNamespaceListerExpansion
}

// vMNodeScrapeNamespaceLister implements the VMNodeScrapeNamespaceLister
// interface.
type vMNodeScrapeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VMNodeScrapes in the indexer for a given namespace.
func (s vMNodeScrapeNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.VMNodeScrape, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VMNodeScrape))
	})
	return ret, err
}

// Get retrieves the VMNodeScrape from the indexer for a given namespace and name.
func (s vMNodeScrapeNamespaceLister) Get(name string) (*v1beta1.VMNodeScrape, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("vmnodescrape"), name)
	}
	return obj.(*v1beta1.VMNodeScrape), nil
}