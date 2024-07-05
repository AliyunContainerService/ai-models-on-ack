# LMDeploy

## Background information

[LMDeploy](https://github.com/InternLM/lmdeploy) LMDeploy is a toolkit for compressing, deploying, and serving LLM, developed by the MMRazor and MMDeploy teams. It has the following core features:  

- <b>Efficient Inference</b>: LMDeploy delivers up to 1.8x higher request throughput than vLLM, by introducing key features like persistent batch(a.k.a. continuous batching), blocked KV cache, dynamic split&fuse, tensor parallelism, high-performance CUDA kernels and so on.
  

- <b>Effective Quantization</b>: LMDeploy supports weight-only and k/v quantization, and the 4-bit inference performance is 2.4x higher than FP16. The quantization quality has been confirmed via OpenCompass evaluation.  
  

- <b>Effortless Distribution Server</b>: Leveraging the request distribution service, LMDeploy facilitates an easy and efficient deployment of multi-model services across multiple machines and cards.
  

- <b>Interactive Inference Mode</b>: By caching the k/v of attention during multi-round dialogue processes, the engine remembers dialogue history, thus avoiding repetitive processing of historical sessions.

## Prerequisites
- An ACK Pro cluster that contains GPU-accelerated nodes is created. The Kubernetes version of the cluster is 1.22 or later. Each GPU-accelerated node provides 16 GB of GPU memory or above. For more information, see [Create an ACK managed cluster](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/create-an-ack-managed-cluster-2).


## QuickStart

### Step 1: Prepare the model data
In this example, the Qwen1.5-4B-Chat model is used to describe how to download a Qwen model, upload the model to Object Storage Service (OSS), and create a persistent volume (PV) and persistent volume claim (PVC) in an ACK cluster.

1. Download the model file.
```bash
yum install git git-lfs
GIT_LFS_SKIP_SMUDGE=1 git clone https://www.modelscope.cn/qwen/Qwen1.5-4B-Chat.git
cd Qwen1.5-4B-Chat
git lfs pull
```

2. Upload the Qwen1.5-4B-Chat model file to OSS.
```bash
ossutil mkdir oss://<Your-Bucket-Name>/Qwen1.5-4B-Chat
ossutil cp -r ./Qwen1.5-4B-Chat oss://<Your-Bucket-Name>/Qwen1.5-4B-Chat
```

3. Configure PVs and PVCs in the destination cluster.
> You need to replace the variables in the file with real values.  
```bash
kubectl apply -f ./yamls/dataset.yaml
```

### Step 2: Deploy an inference service
1. Run the following command to deploy the Qwen1.5-4B-Chat model as an inference service by using LMDeploy:

```bash
kubectl apply -f ./yamls/deploy.yaml
kubectl apply -f ./yamls/service.yaml
```

2. Run the following command to view the details of the inference service:
```bash
$ kubectl get po|grep lmdeploy
---
lmdeploy-6f54847c94-k54kt   1/1     Running   0          27s
```
The output indicates that the inference service is running as expected and is ready to provide services.


### Step 3: Verify the inference service
1. Run the following command to create a port forwarding rule between the inference service and the local environment:
```bash
kubectl port-forward svc/lmdeploy-service 8000:8000
```
Expected output:
```text
Forwarding from 127.0.0.1:8000 -> 8000
Forwarding from [::1]:8000 -> 8000
```
2. Run the following command to send an inference request:
```bash
curl http://localhost:8000/v1/chat/completions \
 -H "Content-Type: application/json" \
 -d '{"model": "qwen", "messages": [{"role": "user", "content": "测试一下"}], "max_tokens": 10, "temperature": 0.7, "top_p": 0.9, "seed": 10}'
```
Expected output:
```text
{"id":"1","object":"chat.completion","created":1720145825,"model":"qwen","choices":[{"index":0,"message":{"role":"assistant","content":"好的，有什么我可以帮助你的吗？"},"logprobs":null,"finish_reason":"stop"}],"usage":{"prompt_tokens":21,"total_tokens":29,"completion_tokens":8}}
```
### (Optional) Step 4: Delete the environment
If you no longer need the resources, delete the resources at the earliest opportunity.
- Run the following command to delete the inference service:
```bash
kubectl delete -f ./yamls
```

## Release Tag

|   tag    |   Date    |  release  |
|:--------:|:---------:|:---------:|
|  v0.4.2  |  2024-07  |   init    |           


