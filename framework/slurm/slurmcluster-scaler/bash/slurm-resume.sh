#!/bin/bash
SCALER=/usr/bin/slurmcluster-scaler
SLURM_NAMESPACE=$(cat /etc/slurmd-podinfo/namespace)
SLURM_CLUSTER=$(grep "^kai.alibabacloud.com/slurm-cluster=" "/etc/slurmd-podinfo/labels" | awk -F "=" '{print $2}' | tr -d \'\")
echo "namespace: $SLURM_NAMESPACE cluster: $SLURM_CLUSTER" >> /var/log/slurm-resume.log
SLURM_NAMESPACE=$SLURM_NAMESPACE SLURM_CLUSTER=$SLURM_CLUSTER KUBERNETES_SERVICE_PORT=$(cat /var/k8s-info/k8s_svc_port) KUBERNETES_SERVICE_HOST=$(cat /var/k8s-info/k8s_svc_host) $SCALER resume "$@" &>> /var/log/slurm-resume.log
