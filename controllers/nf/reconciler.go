/*
Copyright 2023 The Nephio Authors.

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

package nf

import (
	"context"

	hss "github.com/amitinfo2k/sdcore-operator-go/controllers/nf/hss"
	mme "github.com/amitinfo2k/sdcore-operator-go/controllers/nf/mme"
	pcrf "github.com/amitinfo2k/sdcore-operator-go/controllers/nf/pcrf"
	spgwc "github.com/amitinfo2k/sdcore-operator-go/controllers/nf/spgwc"
	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	//	upf "github.com/nephio-project/free5gc/controllers/nf/upf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciles a NFDeployment resource
type NFDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Sets up the controller with the Manager
func (r *NFDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(new(nephiov1alpha1.NFDeployment)).
		Owns(new(appsv1.Deployment)).
		Owns(new(apiv1.ConfigMap)).
		Complete(r)
}

// +kubebuilder:rbac:groups=workload.nephio.org,resources=nfdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=workload.nephio.org,resources=nfdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="ref.nephio.org",resources=configs,verbs=get;list;watch
// +kubebuilder:rbac:groups="k8s.cni.cncf.io",resources=network-attachment-definitions,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NFDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *NFDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("NFDeployment", req.NamespacedName)

	nfDeployment := new(nephiov1alpha1.NFDeployment)
	err := r.Client.Get(ctx, req.NamespacedName, nfDeployment)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("NFDeployment resource not found, ignoring because object must be deleted")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get NFDeployment")
		return reconcile.Result{}, err
	}

	hssReconciler := &hss.HSSDeploymentReconciler{
		Client: r.Client,
		Scheme: r.Scheme,
	}
	pcrfReconciler := &pcrf.PCRFDeploymentReconciler{
		Client: r.Client,
		Scheme: r.Scheme,
	}
	spgwcReconciler := &spgwc.SPGWCDeploymentReconciler{
		Client: r.Client,
		Scheme: r.Scheme,
	}
	mmeReconciler := &mme.MMEDeploymentReconciler{
		Client: r.Client,
		Scheme: r.Scheme,
	}
	config4gReconciler := &mme.Config4GDeploymentReconciler{
		Client: r.Client,
		Scheme: r.Scheme,
	}

	switch nfDeployment.Spec.Provider {
	//	case "upf.free5gc.io":
	//		upfresult, _ := upfReconciler.Reconcile(ctx, req)
	//		return upfresult, nil
	case "hss.sdcore4g.io":
		hssresult, _ := hssReconciler.Reconcile(ctx, req)
		return hssresult, nil
	case "pcrf.sdcore4g.io":
		pcrfresult, _ := pcrfReconciler.Reconcile(ctx, req)
		return pcrfresult, nil
	case "spgwc.sdcore4g.io":
		spgwcresult, _ := spgwcReconciler.Reconcile(ctx, req)
		return spgwcresult, nil
	case "mme.sdcore4g.io":
		mmeresult, _ := mmeReconciler.Reconcile(ctx, req)
		return mmeresult, nil
	case "config4g.sdcore4g.io":
		mmeresult, _ := mmeReconciler.Reconcile(ctx, req)
		return mmeresult, nil
	default:
		log.Info("NFDeployment NOT for SDCore 4G", "nfDeployment.Spec.Provider", nfDeployment.Spec.Provider)
		return reconcile.Result{}, nil
	}
}
