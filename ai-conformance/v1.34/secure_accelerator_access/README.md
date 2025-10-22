MUST: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.



# Tests
# Test 1
**Step 1:** Prepare the test environment, including:
- Creating a Kubernetes 1.34 cluster
- Adding a GPU node pool
- Installing the NVIDIA GPU DRA driver
Followed: https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/scheduling-gpu-using-dra  

**Step 2 [Accessible]:** Deploy a Pod on a node with available accelerator(s), and ensure the container within the Pod explicitly requests accelerator resources.
Inside the running container, execute a command to detect the accelerator device. This command should succeed and output the model of the accelerator device currently used by the container.
```shell
$ kubectl apply -f - <<eof
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: single-gpu
spec:
  spec:
    devices:
      requests:
      - exactly:
          allocationMode: ExactCount
          deviceClassName: gpu.nvidia.com
          count: 1
        name: gpu
eof
```
```shell
$ kubectl apply -f - <<eof
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dra-gpu-success-access
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dra-gpu-success-access
  template:
    metadata:
      labels:
        app: dra-gpu-success-access
    spec:
      containers:
      - name: probe
        image: ubuntu:22.04
        command: ["bash", "-c"]
        args: ["while [ 1 ]; do date; nvidia-smi -L; sleep 60; done"]
        resources:
          claims:
          - name: single-gpu
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: single-gpu
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
eof
```
**Validated: The ResourceClaim is created successfully. The Pod can output accelerator information without any errors.**
```shell
$ kubectl get resourceclaims
NAME                                                       STATE                AGE
dra-gpu-success-access-5b5b546cff-2qftl-single-gpu-z78wc   allocated,reserved   6s

$ kubectl get pods -n nvidia-dra-driver-gpu
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-6658b47869-k49jn   1/1     Running   0          103m
nvidia-dra-driver-gpu-kubelet-plugin-8kz5p          2/2     Running   0          103m
nvidia-dra-driver-gpu-kubelet-plugin-f2dmm          2/2     Running   0          45m

$ kubectl get po -l app=dra-gpu-success-access
NAME                                      READY   STATUS    RESTARTS   AGE
dra-gpu-success-access-5b5b546cff-2qftl   1/1     Running   0          17s

$ kubectl logs -l app=dra-gpu-success-access
Mon Oct 13 08:33:46 UTC 2025
GPU 0: NVIDIA A10 (UUID: GPU-ded74eb2-2ec4-b71d-6465-67f1a532096b)
```
**Step 3 [Prevent unexpected access]:**
Deploy a Pod on a node with available accelerator(s), but ensure the container inside the Pod does not request accelerator resources.
Inside the running container, execute a command to detect the accelerator device. This command should fail or return a message indicating no accelerator device was found.
```shell
$ kubectl apply -f - <<eof
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dra-gpu-fail-access
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dra-gpu-fail-access
  template:
    metadata:
      labels:
        app: dra-gpu-fail-access
    spec:
      containers:
      - name: probe
        image: ubuntu:22.04
        command: ["bash", "-c"]
        args: ["while [ 1 ]; do date; nvidia-smi -L; sleep 60; done"]
        # resources:
        #   claims:
        #   - name: single-gpu
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: single-gpu
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
eof
```
**Validated: The Pod cannot access the accelerator device that was not requested.**
```shell
$ kubectl logs -l app=dra-gpu-fail-access
Mon Oct 13 08:35:27 UTC 2025
bash: line 1: nvidia-smi: command not found
```

## Test 2
**Step 1:** Create two Pods, each is allocated an accelerator resource. Execute a command in one Pod to attempt to access the other Pod’s accelerator, and should be denied.

This can be verified by running this test https://github.com/kubernetes/kubernetes/blob/v1.34.1/test/e2e/dra/dra.go#L180 

With a 1.34 ACK cluster:
```shell
$ make WHAT="github.com/onsi/ginkgo/v2/ginkgo k8s.io/kubernetes/test/e2e/e2e.test" && KUBERNETES_PROVIDER=local hack/ginkgo-e2e.sh -ginkgo.focus='must map configs and devices to the right containers'
+++ [1016 16:58:19] Building go targets for linux/amd64
    github.com/onsi/ginkgo/v2/ginkgo (non-static)
    k8s.io/kubernetes/test/e2e/e2e.test (test)
Setting up for KUBERNETES_PROVIDER="local".
Skeleton Provider: prepare-e2e not implemented
KUBE_MASTER_IP: 
KUBE_MASTER: 
  I1016 16:58:27.649393  466863 e2e.go:109] Starting e2e run "fdbf8455-7688-4398-9635-5ec34d824831" on Ginkgo node 1
Running Suite: Kubernetes e2e suite - /root/test/kubernetes/_output/bin
=======================================================================
Random Seed: 1760605107 - will randomize all specs

Will run 1 of 7137 specs
•

Ran 1 of 7137 Specs in 17.041 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 7136 Skipped
PASS

Ginkgo ran 1 suite in 17.734968102s
Test Suite Passed
```
