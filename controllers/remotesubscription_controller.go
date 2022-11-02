/*
Copyright 2022 Ken Moini.

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

package controllers

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	//"sigs.k8s.io/controller-runtime/pkg/log"

	configv1alpha1 "github.com/kenmoini/ozone-operator/api/v1alpha1"
	pkgSrvr "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
)

// RemoteSubscriptionReconciler reconciles a RemoteSubscription object
type RemoteSubscriptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//===========================================================================================
// RBAC GENERATORS
//===========================================================================================
//+kubebuilder:rbac:groups=config.operator.o3,resources=remotesubscriptions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.operator.o3,resources=remotesubscriptions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.operator.o3,resources=remotesubscriptions/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

//+kubebuilder:rbac:groups=packages.operators.coreos.com,resources=PackageManifest,verbs=list;get

// ===========================================================================================
// RECONCILE
// ===========================================================================================
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RemoteSubscription object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *RemoteSubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)

	controllerLog.Info("starting reconciliation loop")

	currentConfig, _ := config.GetConfig()
	clusterEndpoint := currentConfig.Host

	controllerLog.Info("Hub Cluster: " + clusterEndpoint)

	lPrefix := "rsub[" + req.Name + "] "

	// Get the object that triggered the reconcile request
	controllerLog.Info(lPrefix + "Getting RemoteSubscription object")

	// Fetch the RemoteSubscription instance
	remoteSubscription := &configv1alpha1.RemoteSubscription{}
	err := r.Get(ctx, req.NamespacedName, remoteSubscription)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			controllerLog.Error(err, lPrefix+"RemoteSubscription resource not found on the cluster.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		controllerLog.Error(err, lPrefix+"Failed to get RemoteSubscription")
		return ctrl.Result{}, err
	}

	// Log the used remoteSubscription
	controllerLog.Info(lPrefix + "RemoteSubscription loaded!  Found '" + remoteSubscription.Name + "' in 'namespace/" + remoteSubscription.Namespace + "'")

	// Check for the existence of the Operator PackageManifest on the hub cluster
	targetPackageManifest := &pkgSrvr.PackageManifest{}

	lPrefix = lPrefix + "pm[" + remoteSubscription.Spec.Operator.PackageName + "] "

	controllerLog.Info(lPrefix + "Checking for the existence of the PackageManifest on the hub cluster")

	err = r.Get(ctx, client.ObjectKey{
		Namespace: remoteSubscription.Spec.Operator.PackageNamespace,
		Name:      remoteSubscription.Spec.Operator.PackageName,
	}, targetPackageManifest)

	if err != nil {
		controllerLog.Error(err, lPrefix+"Failed to get PackageManifest/"+remoteSubscription.Spec.Operator.PackageName+" in "+remoteSubscription.Spec.Operator.PackageNamespace)
		return ctrl.Result{}, err
	}

	controllerLog.Info(lPrefix + "PackageManifest found!  Found '" + targetPackageManifest.Name + "' in '" + targetPackageManifest.Namespace + "'")

	// Dump the data from the PackageManifest
	//log.Printf("PackageManifest data: %v", targetPackageManifest)

	pmDefaultChannel := targetPackageManifest.Status.DefaultChannel
	pmTargetChannel := SetDefaultString(pmDefaultChannel, remoteSubscription.Spec.Operator.Channel)

	// Dump debug data
	controllerLog.Info(lPrefix + "PackageManifest default channel: " + pmDefaultChannel)
	controllerLog.Info(lPrefix + "PackageManifest specified channel: " + remoteSubscription.Spec.Operator.Channel)
	controllerLog.Info(lPrefix + "PackageManifest target channel: " + pmTargetChannel)

	// Get the channel from the PackageManifest
	var pmTargetChannelData *pkgSrvr.PackageChannel
	for _, channel := range targetPackageManifest.Status.Channels {
		if channel.Name == pmTargetChannel {
			pmTargetChannelData = &channel
		}
	}
	// Make sure the pmTargetChannelData is not nil
	if pmTargetChannelData == nil {
		controllerLog.Error(err, lPrefix+"Failed to find the target channel in the PackageManifest")
		return ctrl.Result{}, err
	}

	// Dump the channel data from the PackageManifest
	log.Printf("channel data: %v", pmTargetChannelData)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RemoteSubscriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.RemoteSubscription{}).
		Complete(r)
}
