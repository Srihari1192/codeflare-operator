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

package e2e

import (
	"testing"

	. "github.com/onsi/gomega"
	mcadv1beta1 "github.com/project-codeflare/multi-cluster-app-dispatcher/pkg/apis/controller/v1beta1"
	rayv1alpha1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/project-codeflare/codeflare-operator/test/support"
)

// Ray cluster up
func TestMNISTRayClusterUp(t *testing.T) {
	test := With(t)
	test.T().Parallel()

	// Create a namespace
	namespace := test.NewTestNamespace()

	// Test configuration
	config := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mnist-raycluster-sdk",
			Namespace: namespace.Name,
		},
		BinaryData: map[string][]byte{
			// SDK script
			"mnist_raycluster_sdk.py": ReadFile(test, "mnist_raycluster_sdk.py"),
			// pip requirements
			"requirements.txt": ReadFile(test, "mnist_pip_requirements.txt"),
			// MNIST training script
			"mnist.py": ReadFile(test, "mnist.py"),
		},
		Immutable: Ptr(true),
	}
	config, err := test.Client().Core().CoreV1().ConfigMaps(namespace.Name).Create(test.Ctx(), config, metav1.CreateOptions{})
	test.Expect(err).NotTo(HaveOccurred())
	test.T().Logf("Created ConfigMap %s/%s successfully", config.Namespace, config.Name)

	// SDK client RBAC
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdk-user",
			Namespace: namespace.Name,
		},
	}
	serviceAccount, err = test.Client().Core().CoreV1().ServiceAccounts(namespace.Name).Create(test.Ctx(), serviceAccount, metav1.CreateOptions{})
	test.Expect(err).NotTo(HaveOccurred())

	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rbacv1.SchemeGroupVersion.String(),
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdk",
			Namespace: namespace.Name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get", "create", "delete", "list", "patch", "update"},
				APIGroups: []string{mcadv1beta1.GroupName},
				Resources: []string{"appwrappers"},
			},
			{
				Verbs:     []string{"get", "list"},
				APIGroups: []string{rayv1alpha1.GroupVersion.Group},
				Resources: []string{"rayclusters", "rayclusters/status"},
			},
			{
				Verbs:     []string{"get", "list"},
				APIGroups: []string{"route.openshift.io"},
				Resources: []string{"routes"},
			},
		},
	}
	role, err = test.Client().Core().RbacV1().Roles(namespace.Name).Create(test.Ctx(), role, metav1.CreateOptions{})
	test.Expect(err).NotTo(HaveOccurred())

	roleBinding := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rbacv1.SchemeGroupVersion.String(),
			Kind:       "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "Role",
			Name:     role.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				APIGroup:  corev1.SchemeGroupVersion.Group,
				Name:      serviceAccount.Name,
				Namespace: serviceAccount.Namespace,
			},
		},
	}
	_, err = test.Client().Core().RbacV1().RoleBindings(namespace.Name).Create(test.Ctx(), roleBinding, metav1.CreateOptions{})
	test.Expect(err).NotTo(HaveOccurred())
}
