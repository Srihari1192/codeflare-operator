package e2e

import (
	"testing"

	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/project-codeflare/codeflare-operator/test/support"
)

func TestMnistJobSubmit(t *testing.T) {

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
			"mnist_rayjob.py": ReadFile(test, "mnist_rayjob.py"),
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

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: batchv1.SchemeGroupVersion.String(),
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdk",
			Namespace: namespace.Name,
		},
		Spec: batchv1.JobSpec{
			Completions:  Ptr(int32(1)),
			Parallelism:  Ptr(int32(1)),
			BackoffLimit: Ptr(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
							// FIXME: switch to base Python image once the dependency on OpenShift CLI is removed
							// See https://github.com/project-codeflare/codeflare-sdk/pull/146
							Image:   "quay.io/opendatahub/notebooks:jupyter-minimal-ubi8-python-3.8-4c8f26e",
							Command: []string{"/bin/sh", "-c", "pip install codeflare-sdk==" + GetCodeFlareSDKVersion() + " && cp /test/* . && python mnist_raycluster_sdk.py" + " " + namespace.Name},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "test",
									MountPath: "/test",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "test",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: config.Name,
									},
								},
							},
						},
					},
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: serviceAccount.Name,
				},
			},
		},
	}
	job, err = test.Client().Core().BatchV1().Jobs(namespace.Name).Create(test.Ctx(), job, metav1.CreateOptions{})
	test.Expect(err).NotTo(HaveOccurred())
	test.T().Logf("Created Job %s/%s successfully", job.Namespace, job.Name)

	// Retrieving the job logs once it has completed or timed out
	defer WriteJobLogs(test, job.Namespace, job.Name)

	test.T().Logf("Waiting for Job %s/%s to complete", job.Namespace, job.Name)
	test.Eventually(Job(test, job.Namespace, job.Name), TestTimeoutLong).Should(
		Or(
			WithTransform(ConditionStatus(batchv1.JobComplete), Equal(corev1.ConditionTrue)),
			WithTransform(ConditionStatus(batchv1.JobFailed), Equal(corev1.ConditionTrue)),
		))

	// Assert the job has completed successfully
	test.Expect(GetJob(test, job.Namespace, job.Name)).
		To(WithTransform(ConditionStatus(batchv1.JobComplete), Equal(corev1.ConditionTrue)))
}
