/*
Copyright 2023 Frédéric Boltz.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	nodemanagerv1alpha1 "github.com/Fred78290/kubernetes-desktop-autoscaler/pkg/apis/nodemanager/v1alpha1"
	versioned "github.com/Fred78290/kubernetes-desktop-autoscaler/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/Fred78290/kubernetes-desktop-autoscaler/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/Fred78290/kubernetes-desktop-autoscaler/pkg/generated/listers/nodemanager/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ManagedNodeInformer provides access to a shared informer and lister for
// ManagedNodes.
type ManagedNodeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ManagedNodeLister
}

type managedNodeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewManagedNodeInformer constructs a new informer for ManagedNode type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewManagedNodeInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredManagedNodeInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredManagedNodeInformer constructs a new informer for ManagedNode type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredManagedNodeInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NodemanagerV1alpha1().ManagedNodes().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NodemanagerV1alpha1().ManagedNodes().Watch(context.TODO(), options)
			},
		},
		&nodemanagerv1alpha1.ManagedNode{},
		resyncPeriod,
		indexers,
	)
}

func (f *managedNodeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredManagedNodeInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *managedNodeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&nodemanagerv1alpha1.ManagedNode{}, f.defaultInformer)
}

func (f *managedNodeInformer) Lister() v1alpha1.ManagedNodeLister {
	return v1alpha1.NewManagedNodeLister(f.Informer().GetIndexer())
}
