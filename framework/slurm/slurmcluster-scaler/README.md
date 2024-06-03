# SlurmCluster Scaler

SlurmCluster Scaler is the tool which helps to scale slurm worker nodes
when using ack-slurm-operator on ACK.

## How to use

1. Build and put the binary in your image which contains slurmctld by adding
these lines in dockerfile.
``` bash
# Build scaler
FROM golang:1.22 as scalerBuilder
WORKDIR /usr/local/go/src/github.com/AliyunContainerService/
RUN apt-get update && apt install git
RUN git clone https://github.com/AliyunContainerService/ai-models-on-ack.git
WORKDIR /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o slurmctldcluster-scaler main.go

# Begining of your dockerfile
# End of your dockerfile

# Copy scaler and scripts
COPY --from=0 /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/slurmctld-copilot /usr/bin/slurmctld-copilot 
COPY --from=0 /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/bash/slurm-resume.sh /usr/bin/slurm-resume
COPY --from=0 /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/bash/slurm-suspend.sh /usr/bin/slurm-suspend
RUN chmod +x /usr/bin/slurm-suspend /usr/bin/slurm-resume /usr/bin/slurmctld-copilot 
```

2. Add role for slurmcluster.
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Release.Name }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Release.Name }}
rules:
- apiGroups: ["kai.alibabacloud.com"]
  resources: ["slurmclusters"]
  verbs: ["get", "watch", "list", "update", "patch"]
  resourceNames: ["{{ .Release.Name }}"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Release.Name }}
subjects:
- kind: ServiceAccount
  name: {{ .Release.Name }}
roleRef:
  kind: Role
  name: {{ .Release.Name }}
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: kai.alibabacloud.com/v1
kind: SlurmCluster
...
spec:
  headGroupSpec:
    template:
      spec:
        serviceAccountName: {{ .Release.Name }}
...
```

3. Add configuration in slurm.conf.
```
SuspendTimeout=600
ResumeTimeout=600
SuspendTime=600
ResumeRate=1
SuspendRate=1
CommunicationParameters=NoAddrCache
ReconfigFlags=KeepPowerSaveSettings
# NodeName must be in format: ＄{cluster_name}-worker-${group_name}-
NodeName=slurm-job-demo-worker-cpu-[0-10] Feature=cloud State=CLOUD
SuspendProgram="/usr/bin/slurm-suspend"
ResumeProgram="/usr/bin/slurm-resume"
```