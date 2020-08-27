/*
Copyright The Kubernetes Authors.

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
	v1alpha1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/apis/ako/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// HTTPRuleLister helps list HTTPRules.
type HTTPRuleLister interface {
	// List lists all HTTPRules in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.HTTPRule, err error)
	// HTTPRules returns an object that can list and get HTTPRules.
	HTTPRules(namespace string) HTTPRuleNamespaceLister
	HTTPRuleListerExpansion
}

// hTTPRuleLister implements the HTTPRuleLister interface.
type hTTPRuleLister struct {
	indexer cache.Indexer
}

// NewHTTPRuleLister returns a new HTTPRuleLister.
func NewHTTPRuleLister(indexer cache.Indexer) HTTPRuleLister {
	return &hTTPRuleLister{indexer: indexer}
}

// List lists all HTTPRules in the indexer.
func (s *hTTPRuleLister) List(selector labels.Selector) (ret []*v1alpha1.HTTPRule, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.HTTPRule))
	})
	return ret, err
}

// HTTPRules returns an object that can list and get HTTPRules.
func (s *hTTPRuleLister) HTTPRules(namespace string) HTTPRuleNamespaceLister {
	return hTTPRuleNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// HTTPRuleNamespaceLister helps list and get HTTPRules.
type HTTPRuleNamespaceLister interface {
	// List lists all HTTPRules in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.HTTPRule, err error)
	// Get retrieves the HTTPRule from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.HTTPRule, error)
	HTTPRuleNamespaceListerExpansion
}

// hTTPRuleNamespaceLister implements the HTTPRuleNamespaceLister
// interface.
type hTTPRuleNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all HTTPRules in the indexer for a given namespace.
func (s hTTPRuleNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.HTTPRule, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.HTTPRule))
	})
	return ret, err
}

// Get retrieves the HTTPRule from the indexer for a given namespace and name.
func (s hTTPRuleNamespaceLister) Get(name string) (*v1alpha1.HTTPRule, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("httprule"), name)
	}
	return obj.(*v1alpha1.HTTPRule), nil
}
