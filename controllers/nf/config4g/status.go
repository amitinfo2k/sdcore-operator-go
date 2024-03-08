package config4g

import (
	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createNfDeploymentStatus(deployment *appsv1.Deployment, config4gDeployment *nephiov1alpha1.NFDeployment) (nephiov1alpha1.NFDeploymentStatus, bool) {
	nfDeploymentStatus := nephiov1alpha1.NFDeploymentStatus{
		ObservedGeneration: int32(deployment.Generation),
		Conditions:         config4gDeployment.Status.Conditions,
	}

	// Return initial status if there are no status update happened for the Config4Gdeployment
	if len(config4gDeployment.Status.Conditions) == 0 {
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Reconciling),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "Config4GDeployment pod(s) is(are) starting.",
			LastTransitionTime: metav1.Now(),
		})

		return nfDeploymentStatus, true
	} else if (len(deployment.Status.Conditions) == 0) && (len(config4gDeployment.Status.Conditions) > 0) {
		return nfDeploymentStatus, false
	}

	// Check the last underlying Deployment status and deduce condition from it
	lastDeploymentCondition := deployment.Status.Conditions[0]
	lastAmfDeploymentCondition := config4gDeployment.Status.Conditions[len(config4gDeployment.Status.Conditions)-1]

	// Deployemnt and Config4GDeployment have different names for processing state, hence we check if one is processing another is reconciling, then state is equal
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
			Message:            "Config4GDeployment pods are available.",
			LastTransitionTime: metav1.Now(),
		})

	case appsv1.DeploymentProgressing:
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Reconciling),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "Config4GDeployment pod(s) is(are) starting.",
			LastTransitionTime: metav1.Now(),
		})

	case appsv1.DeploymentReplicaFailure:
		nfDeploymentStatus.Conditions = append(nfDeploymentStatus.Conditions, metav1.Condition{
			Type:               string(nephiov1alpha1.Stalled),
			Status:             metav1.ConditionFalse,
			Reason:             "MinimumReplicasNotAvailable",
			Message:            "Config4GDeployment pod(s) is(are) failing.",
			LastTransitionTime: metav1.Now(),
		})
	}

	return nfDeploymentStatus, true
}
