package deploymanager

import (
	"time"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilwait "k8s.io/apimachinery/pkg/util/wait"
)

// CreateNamespace creates a namespace in the cluster, ignoring if it already exists
func (t *DeployManager) CreateNamespace(namespace string) error {
	ns := &k8sv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := t.k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// DeleteNamespaceAndWait deletes a namespace and waits on it to terminate
func (t *DeployManager) DeleteNamespaceAndWait(namespace string) error {
	err := t.k8sClient.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	timeout := 200 * time.Second
	interval := 10 * time.Second

	// Wait for namespace to terminate
	err = utilwait.PollImmediate(interval, timeout, func() (done bool, err error) {
		_, err = t.k8sClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
		if !errors.IsNotFound(err) {
			return false, nil
		}
		return true, nil
	})

	return err
}
