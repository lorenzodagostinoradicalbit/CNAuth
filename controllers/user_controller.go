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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lorenzodagostinoradicalbit/CNAuth/api/v1alpha1"
	keysv1alpha1 "github.com/lorenzodagostinoradicalbit/CNAuth/api/v1alpha1"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=keys.cnauth,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keys.cnauth,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=keys.cnauth,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	user, err := r.getUser(ctx, req)
	if err != nil {
		logger.Error(err, "unable to fetch User")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	// Fetch the JWTKey
	RequestedKey := user.Spec.KeyRef
	jwtKeyObj := &keysv1alpha1.JWTKey{}
	jwtNsName := types.NamespacedName{
		Namespace: req.NamespacedName.Namespace,
		Name:      RequestedKey,
	}
	err = r.Client.Get(ctx, jwtNsName, jwtKeyObj)
	if err != nil {
		logger.Error(err, "Unable to fetch jwt object")
		return ctrl.Result{}, err
	}
	logger.Info("Generating JWT token")
	jwtKey := jwtKeyObj.Status.Key
	tokenStr, err := generateJWT(jwtKey, user)
	if err != nil {
		return ctrl.Result{}, err
	}
	user.Status.Token = tokenStr
	r.Status().Update(ctx, user)

	return ctrl.Result{}, nil
}

func generateJWT(key string, user *v1alpha1.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = user.Spec.Name
	// claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (r *UserReconciler) getUser(ctx context.Context, req ctrl.Request) (*keysv1alpha1.User, error) {
	user := &keysv1alpha1.User{}
	if err := r.Client.Get(ctx, req.NamespacedName, user); err != nil {
		return nil, err
	}
	return user, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&keysv1alpha1.User{}).
		Complete(r)
}
