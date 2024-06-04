package common

import (
	"github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/generated/kai.alibabacloud.com/clientset/versioned"
	ctrl "sigs.k8s.io/controller-runtime"
)

func InitClientSet() *versioned.Clientset {
	config := ctrl.GetConfigOrDie()
	clientSet := versioned.NewForConfigOrDie(config)
	return clientSet
}
