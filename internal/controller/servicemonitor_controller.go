/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/ryanwuer/k8s-housekeeper-operator/api/v1alpha1"
)

const (
	finalizerName = "monitoring.cluster.local/finalizer"
)

// ServiceMonitorReconciler reconciles a ServiceMonitor object
type ServiceMonitorReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=monitoring.cluster.local,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.cluster.local,resources=servicemonitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.cluster.local,resources=servicemonitors/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ServiceMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ServiceMonitor instance
	serviceMonitor := &monitoringv1alpha1.ServiceMonitor{}
	err := r.Get(ctx, req.NamespacedName, serviceMonitor)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle deletion
	if !serviceMonitor.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(serviceMonitor, finalizerName) {
			// Remove finalizer
			controllerutil.RemoveFinalizer(serviceMonitor, finalizerName)
			if err := r.Update(ctx, serviceMonitor); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(serviceMonitor, finalizerName) {
		controllerutil.AddFinalizer(serviceMonitor, finalizerName)
		if err := r.Update(ctx, serviceMonitor); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if it's time for the next check
	if serviceMonitor.Status.LastCheckTime != nil {
		nextCheck := serviceMonitor.Status.LastCheckTime.Add(time.Duration(serviceMonitor.Spec.CheckInterval) * time.Second)
		if time.Now().Before(nextCheck) {
			return ctrl.Result{RequeueAfter: time.Until(nextCheck)}, nil
		}
	}

	// List all services
	serviceList := &corev1.ServiceList{}
	if err := r.List(ctx, serviceList); err != nil {
		log.Error(err, "Failed to list services")
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	// Count services with matching ExternalName
	count := int32(0)
	for _, svc := range serviceList.Items {
		if svc.Spec.Type == corev1.ServiceTypeExternalName && svc.Spec.ExternalName == serviceMonitor.Spec.TargetDomain {
			count++
		}
	}

	// Update status
	now := metav1.Now()
	serviceMonitor.Status.LastCheckTime = &now
	serviceMonitor.Status.Count = count

	// Update conditions
	meta.SetStatusCondition(&serviceMonitor.Status.Conditions, metav1.Condition{
		Type:               "Available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "Monitoring",
		Message:            fmt.Sprintf("Found %d services with matching ExternalName", count),
	})

	if err := r.Status().Update(ctx, serviceMonitor); err != nil {
		log.Error(err, "Failed to update ServiceMonitor status")
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	// Schedule next check
	return ctrl.Result{RequeueAfter: time.Duration(serviceMonitor.Spec.CheckInterval) * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ServiceMonitor{}).
		Watches(
			&corev1.Service{},
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}
