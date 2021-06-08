package util

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sHandler struct {
	k8sClient *kubernetes.Clientset
	namespace string
}

func InitSimpleScheduler() (*K8sHandler, error) {
	var config *restclient.Config
	var err error

	if viper.GetBool("InCluster") {
		config, err = restclient.InClusterConfig()
	} else {
		config, err = createOutOfClusterConfig()
	}
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	namespace := viper.GetString("K8sNamespace")

	scheduler := K8sHandler{
		k8sClient: clientset,
		namespace: namespace,
	}

	return &scheduler, nil
}

//createOutOfClusterConfig tries to create the required k8s configuration from well known config paths
func createOutOfClusterConfig() (*restclient.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config, err
}

func (handler *K8sHandler) DeleteK8sJob(jobID string) error {
	k8sJobName := fmt.Sprintf("%s%s", "bakta-job-", jobID)

	delProp := metav1.DeletePropagationForeground

	err := handler.k8sClient.BatchV1().Jobs(handler.namespace).Delete(context.TODO(), k8sJobName, metav1.DeleteOptions{
		PropagationPolicy: &delProp,
	})
	if err != nil && !errors.IsNotFound(err) {
		glog.Errorln(err.Error())
		return err
	}

	//TODO: Only in debug mode
	if errors.IsNotFound(err) {
		glog.Errorln(err.Error())
	}

	return nil
}
