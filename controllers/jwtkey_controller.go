/*
Copyright 2023.

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	keysv1alpha1 "github.com/lorenzodagostinoradicalbit/CNAuth/api/v1alpha1"
)

// JWTKeyReconciler reconciles a JWTKey object
type JWTKeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=keys.cnauth,resources=jwtkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keys.cnauth,resources=jwtkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=keys.cnauth,resources=jwtkeys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JWTKey object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *JWTKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	key, err := r.getJWTKey(ctx, req)
	if err != nil {
		logger.Error(err, "unable to fetch JWTKey")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	// Simply set the status key as the Spec key
	key.Status.Key = key.Spec.Key
	r.Status().Update(ctx, key)

	return ctrl.Result{}, nil
}

func (r *JWTKeyReconciler) getJWTKey(ctx context.Context, req ctrl.Request) (*keysv1alpha1.JWTKey, error) {
	key := &keysv1alpha1.JWTKey{}
	if err := r.Client.Get(ctx, req.NamespacedName, key); err != nil {
		return nil, err
	}
	return key, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JWTKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&keysv1alpha1.JWTKey{}).
		Complete(r)
}
