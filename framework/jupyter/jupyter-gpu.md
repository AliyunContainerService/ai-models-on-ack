# Jupyter GPU

## QuickStart

### Prerequisites

- An ASK cluster is created. For more information,
  see [Create an ASK cluster](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/user-guide/create-an-ask-cluster-2?spm=a2c63.p38356.0.0.664265cdTbNZo1#task-e3c-311-ydb)
- The cluster runs as expected. You can log on to the Container Service for Kubernetes (ACK) console, navigate to the
  Clusters page, and then check whether the cluster is in the Running state.

### Deploy Stable Diffusion

1. create jupyter deployment & service

```bash
kubectl apply -f jupyter-gpu.yaml
```  

2. wait deployment ready

```bash
kubectl get po |grep notebook

# NAME                       READY   STATUS    RESTARTS   AGE
# notebook-d68d854c9-ptvtp   1/1     Running   0          8m5s
```

3. connect to the Jupyter Notebook

Run the following command to port-forward:

```
kubectl port-forward -n <namespace> service/notebook-svc 8888:8888
```

And then open the console using the following URL:

```
http://localhost:8888
```

![jupyter-gpu](jupyter-gpu.jpg "jupyter-gpu")

Run the following command to check the gpu device:

```python
! nvidia-smi
```