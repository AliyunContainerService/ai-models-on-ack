# FastChat

---
## QuickStart  

### Prerequisites
- An ASK cluster is created. For more information, see [Create an ASK cluster](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/user-guide/create-an-ask-cluster-2?spm=a2c63.p38356.0.0.664265cdTbNZo1#task-e3c-311-ydb)
- The cluster runs as expected. You can log on to the Container Service for Kubernetes (ACK) console, navigate to the Clusters page, and then check whether the cluster is in the Running state.

### Deploy Stable Diffusion
```bash
# 1. create deployment
# only for gpu
fastchat-deployment.yaml

# 2. create loadbalancer service
kubectl apply -f fastchat-svc.yaml

# 3. wait deployment ready

# 4. get ExternalIP
kubectl get svc fastchat-svc

# 5. open http://${ExternalIP}:7860 
```
