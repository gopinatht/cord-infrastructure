package main

import (
	"encoding/json"
	"log"
	"time"

	"k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func watchReplicaSet(c *config, clientset *kubernetes.Clientset, restClient rest.Interface) {
	watchlist := cache.NewListWatchFromClient(restClient, "replicasets", corev1.NamespaceAll, fields.Everything())

	// Wrap the returned watchlist to workaround the inability to include
	// the `IncludeUninitialized` list option when setting up watch clients.
	includeUninitializedWatchlist := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist.Watch(options)
		},
	}

	resyncPeriod := 30 * time.Second

	_, controller := cache.NewInformer(includeUninitializedWatchlist, &v1beta2.ReplicaSet{}, resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				err := initializeReplicaset(obj.(*v1beta2.ReplicaSet), c, clientset)
				if err != nil {
					log.Println(err)
				}
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)
}

func initializeReplicaset(replicaset *v1beta2.ReplicaSet, c *config, clientset *kubernetes.Clientset) error {
	if replicaset.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := replicaset.ObjectMeta.GetInitializers().Pending

		if initializerName == pendingInitializers[0].Name {
			log.Printf("Initializing replicaset: %s", replicaset.Name)

			// o, err := runtime.NewScheme().DeepCopy(replicaset)
			// if err != nil {
			// 	return err
			// }
			// initializedReplicaSet := o.(*v1beta2.ReplicaSet)

			initializedReplicaSet := replicaset.DeepCopy()

			if requireAnnotation {
				a := replicaset.ObjectMeta.GetAnnotations()
				_, ok := a[annotation]
				if !ok {
					log.Printf("Required '%s' annotation missing; skipping cord infra container injection", annotation)
					_, err := clientset.AppsV1beta2().ReplicaSets(replicaset.Namespace).Update(initializedReplicaSet)
					if err != nil {
						return err
					}
					return nil
				}
			}

			// Modify the ReplicaSet's Pod template to include the cord infra container
			// and configuration volume. Then patch the original replicaset.
			initializedReplicaSet.Spec.Template.Spec.Containers = append(replicaset.Spec.Template.Spec.Containers, c.Containers...)
			initializedReplicaSet.Spec.Template.Spec.Volumes = append(replicaset.Spec.Template.Spec.Volumes, c.Volumes...)

			oldData, err := json.Marshal(replicaset)
			if err != nil {
				return err
			}

			newData, err := json.Marshal(initializedReplicaSet)
			if err != nil {
				return err
			}

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1beta2.ReplicaSet{})
			if err != nil {
				return err
			}

			_, err = clientset.AppsV1beta2().ReplicaSets(replicaset.Namespace).Patch(replicaset.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
