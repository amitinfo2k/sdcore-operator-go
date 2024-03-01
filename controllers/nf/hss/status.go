package hss

import (
	"github.com/amitinfo2k/sdcore-operator-go/api/v1alpha1"
	//"github.com/amitinfo2k/sdcore-operator-go/controllers"
	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createNfDeploymentStatus(deployment *appsv1.Deployment, hssDeployment *v1alpha1.HSSDeployment) (nephiov1alpha1.NFDeploymentStatus, bool) {
	nfDeploymentStatus := nephiov1alpha1.NFDeploymentStatus{
		ObservedGeneration: int32(deployment.Generation),
		Conditions:         hssDeployment.Status.Conditions,
	}

	// Return initial status if there are no status update happened for the HSSdeployment
	if len(hssDeployment.Status.Conditions) == 0 {
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Reconciling),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "HSSDeployment pod(s) is(are) starting.",
			LastTransitionTime: metav1.Now(),
		})

		return nfDeploymentStatus, true
	} else if (len(deployment.Status.Conditions) == 0) && (len(hssDeployment.Status.Conditions) > 0) {
		return nfDeploymentStatus, false
	}

	// Check the last underlying Deployment status and deduce condition from it
	lastDeploymentCondition := deployment.Status.Conditions[0]
	lastAmfDeploymentCondition := hssDeployment.Status.Conditions[len(hssDeployment.Status.Conditions)-1]

	// Deployemnt and HSSDeployment have different names for processing state, hence we check if one is processing another is reconciling, then state is equal
	if (lastDeploymentCondition.Type == appsv1.DeploymentProgressing) && (lastAmfDeploymentCondition.Type == string(nephiov1alpha1.Reconciling)) {
		return nfDeploymentStatus, false
	}

	// if both status types are Available, don't update.
	if string(lastDeploymentCondition.Type) == string(lastAmfDeploymentCondition.Type) {
		return nfDeploymentStatus, false
	}

	switch lastDeploymentCondition.Type {
	case appsv1.DeploymentAvailable:
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Available),
			Status:             metav1.ConditionTrue,
			Reason:             "MinimumReplicasAvailable",
			Message:            "HSSDeployment pods are available.",
			LastTransitionTime: metav1.Now(),
		})

	case appsv1.DeploymentProgressing:
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Reconciling),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "HSSDeployment pod(s) is(are) starting.",
			LastTransitionTime: metav1.Now(),
		})

	case appsv1.DeploymentReplicaFailure:
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Stalled),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "HSSDeployment pod(s) is(are) failing.",
			LastTransitionTime: metav1.Now(),
		})
	}

	return nfDeploymentStatus, true
}
