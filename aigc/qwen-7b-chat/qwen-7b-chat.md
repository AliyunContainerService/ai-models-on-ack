# Qwen-7B-Chat

## QuickStart

### Prerequisites

- An ASK cluster is created. For more information,
  see [Create an ASK cluster](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/user-guide/create-an-ask-cluster-2?spm=a2c63.p38356.0.0.664265cdTbNZo1#task-e3c-311-ydb)
- The cluster runs as expected. You can log on to the Container Service for Kubernetes (ACK) console, navigate to the
  Clusters page, and then check whether the cluster is in the Running state.

### Deploy Qwen-7B-Chat

```bash
# 1. create deployment & service
# for cpu
kubectl apply -f qwen-7b-chat-cpu.yaml
# for gpu
kubectl apply -f qwen-7b-chat-gpu.yaml

# 3. wait deployment ready

# 4. get ExternalIP
kubectl get svc qwen-7b-chat-svc

# 5. open http://${ExternalIP}:8000
```
