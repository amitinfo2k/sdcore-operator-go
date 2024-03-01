package hss

import (
	"context"
	"fmt"
	"time"

	"github.com/amitinfo2k/sdcore-operator-go/api/v1alpha1"
	"github.com/amitinfo2k/sdcore-operator-go/controllers"

	//
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciles a HSSDeployment resource
type HSSDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Sets up the controller with the Manager
func (r *HSSDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(new(v1alpha1.HSSDeployment)).
		Owns(new(appsv1.Deployment)).
		Owns(new(apiv1.ConfigMap)).
		Complete(r)
}

// +kubebuilder:rbac:groups=workload.nephio.org,resources=hssdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=workload.nephio.org,resources=hssdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="k8s.cni.cncf.io",resources=network-attachment-definitions,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HSSDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *HSSDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("HSSDeployment", req.NamespacedName)

	hssDeployment := new(v1alpha1.HSSDeployment)
	err := r.Client.Get(ctx, req.NamespacedName, hssDeployment)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("HSSDeployment resource not found, ignoring sibecausence object must be deleted")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get HSSDeployment")
		return reconcile.Result{}, err
	}

	namespace := hssDeployment.Namespace

	configMapFound := false
	scriptConfigMapFound := false
	configMapName := "hss-configs"
	ScriptConfigMapName := "hss-scripts"
	log.Info("Reconcile++ ", "configMapName = ", configMapName)
	var configMapVersion string
	var scriptConfMapVersion string
	currentConfigMap := new(apiv1.ConfigMap)
	if err := r.Client.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: namespace}, currentConfigMap); err == nil {
		configMapFound = true
		configMapVersion = currentConfigMap.ResourceVersion
	}
	log.Info("Reconcile", "configMapFound=", configMapFound, ",configMapVersion:", configMapVersion)
	currentConfigMap = new(apiv1.ConfigMap)

	if err := r.Client.Get(ctx, types.NamespacedName{Name: ScriptConfigMapName, Namespace: namespace}, currentConfigMap); err == nil {
		scriptConfigMapFound = true
		scriptConfMapVersion = currentConfigMap.ResourceVersion
	}
	log.Info("Reconcile", "scriptConfigMapFound=", scriptConfigMapFound, ",configMapVersion:", scriptConfMapVersion)

	serviceFound := false
	serviceName := hssDeployment.Name
	currentService := new(apiv1.Service)
	if err := r.Client.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: namespace}, currentService); err == nil {
		serviceFound = true
	}

	deploymentFound := false
	deploymentName := hssDeployment.Name
	currentDeployment := new(appsv1.Deployment)
	if err := r.Client.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: namespace}, currentDeployment); err == nil {
		deploymentFound = true
	}

	if deploymentFound {
		deployment := currentDeployment.DeepCopy()

		// Updating HSSDeployment status. On the first sets the first Condition to Reconciling.
		// On the subsequent runs it gets undelying depoyment Conditions and use the last one to decide if status has to be updated.
		if deployment.DeletionTimestamp == nil {
			if err := r.syncStatus(ctx, deployment, hssDeployment); err != nil {
				log.Error(err, "Failed to update status")
				return reconcile.Result{}, err
			}
		}

		if currentDeployment.Spec.Template.Annotations[controllers.ConfigMapVersionAnnotation] != scriptConfMapVersion {
			log.Info("ConfigMap has been updated, rolling Deployment pods", "Deployment.namespace", currentDeployment.Namespace, "Deployment.name", currentDeployment.Name)
			log.Info("Reconcile", "configMapVersion:", configMapVersion, ",Annotations:", currentDeployment.Spec.Template.Annotations[controllers.ConfigMapVersionAnnotation])
			currentDeployment.Spec.Template.Annotations[controllers.ConfigMapVersionAnnotation] = configMapVersion

			if err := r.Update(ctx, currentDeployment); err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.namespace", currentDeployment.Namespace, "Deployment.name", currentDeployment.Name)
				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	if configMap, err := createConfigMap(log, hssDeployment); err == nil {
		if !configMapFound {
			log.Info("Creating ConfigMap", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)

			// Set the controller reference, specifying that HSSDeployment controling underlying deployment
			if err := ctrl.SetControllerReference(hssDeployment, configMap, r.Scheme); err != nil {
				log.Error(err, "Got error while setting Owner reference on configmap.", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)
			}

			if err := r.Client.Create(ctx, configMap); err != nil {
				log.Error(err, "Failed to create ConfigMap", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)
				return reconcile.Result{}, err
			}

			configMapVersion = configMap.ResourceVersion
		}
	} else {
		log.Error(err, "Failed to create ConfigMap")
		return reconcile.Result{}, err
	} //

	if configMap, err := createScriptConfigMap(log, hssDeployment); err == nil {
		if !scriptConfigMapFound {
			log.Info("Creating ConfigMap", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)

			// Set the controller reference, specifying that HSSDeployment controling underlying deployment
			if err := ctrl.SetControllerReference(hssDeployment, configMap, r.Scheme); err != nil {
				log.Error(err, "Got error while setting Owner reference on configmap.", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)
			}

			if err := r.Client.Create(ctx, configMap); err != nil {
				log.Error(err, "Failed to create ConfigMap", "ConfigMap.namespace", configMap.Namespace, "ConfigMap.name", configMap.Name)
				return reconcile.Result{}, err
			}

			configMapVersion = configMap.ResourceVersion
		}
	} else {
		log.Error(err, "Failed to create Script ConfigMap")
		return reconcile.Result{}, err
	}

	if !serviceFound {
		service := createService(hssDeployment)

		log.Info("Creating HSSDeployment service", "Service.namespace", service.Namespace, "Service.name", service.Name)

		// Set the controller reference, specifying that HSSDeployment controling underlying deployment
		if err := ctrl.SetControllerReference(hssDeployment, service, r.Scheme); err != nil {
			log.Error(err, "Got error while setting Owner reference on HSS service", "Service.namespace", service.Namespace, "Service.name", service.Name)
		}

		if err := r.Client.Create(ctx, service); err != nil {
			log.Error(err, "Failed to create Service", "Service.namespace", service.Namespace, "Service.name", service.Name)
			return reconcile.Result{}, err
		}
	}

	if deployment, err := createDeployment(log, configMapVersion, hssDeployment); err == nil {
		//if deployment, err := createDeployment(log, configMapVersion, hssDeployment); err == nil {
		if !deploymentFound {
			// Only create Deployment in case all required NADs are present. Otherwise Requeue in 10 sec.
			//if ok := controllers.ValidateNetworkAttachmentDefinitions(ctx, r.Client, log, hssDeployment.Kind, deployment); ok {
			// Set the controller reference, specifying that HSSDeployment controls the underlying Deployment
			if err := ctrl.SetControllerReference(hssDeployment, deployment, r.Scheme); err != nil {
				log.Error(err, "Got error while setting Owner reference on deployment", "Deployment.namespace", deployment.Namespace, "Deployment.name", deployment.Name)
			}

			log.Info("Creating Deployment", "Deployment.namespace", deployment.Namespace, "Deployment.name", deployment.Name)
			if err := r.Client.Create(ctx, deployment); err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.namespace", deployment.Namespace, "Deployment.name", deployment.Name)
			}

			// TODO(tliron): explain why we need requeueing (do we?)
			//return reconcile.Result{RequeueAfter: time.Duration(30) * time.Second}, nil
			//} else {
			//		log.Info("Not all NetworkAttachDefinitions available in current namespace, requeuing")
			return reconcile.Result{RequeueAfter: time.Duration(10) * time.Second}, nil
			//	}
		}
	} else {
		log.Error(err, fmt.Sprintf("Failed to create Deployment %s\n", err.Error()))
		return reconcile.Result{}, err
	}
	log.Info("Reconcile--")
	return reconcile.Result{}, nil
}

func (r *HSSDeploymentReconciler) syncStatus(ctx context.Context, deployment *appsv1.Deployment, hssDeployment *v1alpha1.HSSDeployment) error {
	if nfDeploymentStatus, update := createNfDeploymentStatus(deployment, hssDeployment); update {
		hssDeployment = hssDeployment.DeepCopy()
		hssDeployment.Status.NFDeploymentStatus = nfDeploymentStatus
		return r.Status().Update(ctx, hssDeployment)
	} else {
		return nil
	}
}
