# FastChat

## QuickStart

### Prerequisites

- An ASK cluster is created. For more information,
  see [Create an ASK cluster](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/user-guide/create-an-ask-cluster-2?spm=a2c63.p38356.0.0.664265cdTbNZo1#task-e3c-311-ydb)
- The cluster runs as expected. You can log on to the Container Service for Kubernetes (ACK) console, navigate to the
  Clusters page, and then check whether the cluster is in the Running state.

### FastChat

[FastChat](https://github.com/lm-sys/FastChat) is an open platform for training, serving, and evaluating large language
model based chatbots.

### Deploy FastChat

1. create deployment & service

```bash
# only for gpu
kubectl apply -f fastchat-gpu.yaml
```

2. wait deployment ready

```bash
$ kubectl get po|grep fastchat
---
NAME                                READY   STATUS    RESTARTS   AGE
fastchat-65f7cbfbc5-gb7wd           1/1     Running   0          30m
```

3. get ExternalIP

```bash
kubectl get svc fastchat-svc
```

4. open http://${ExternalIP}:7860
   ![fastchat](fastchat.jpg "fastchat")
