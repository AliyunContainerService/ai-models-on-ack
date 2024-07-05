# LMDeploy

## Background information

[vLLM](https://github.com/vllm-project/vllm) vLLM is a high-performance and easy-to-use LLM inference service framework.
vLLM supports most commonly used LLMs, including Qwen models. vLLM is powered by technologies such as PagedAttention
optimization, continuous batching, and model quantification to greatly improve the inference efficiency of LLMs. For
more information about the vLLM framework, see [vLLM GitHub repository](https://github.com/vllm-project/vllm).

## Prerequisites

- An ACK Pro cluster that contains GPU-accelerated nodes is created. The Kubernetes version of the cluster is 1.22 or
  later. Each GPU-accelerated node provides 16 GB of GPU memory or above. For more information,
  see [Create an ACK managed cluster](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/create-an-ack-managed-cluster-2).

## QuickStart

### Step 1: Prepare the model data

In this example, the Qwen1.5-4B-Chat model is used to describe how to download a Qwen model, upload the model to Object
Storage Service (OSS), and create a persistent volume (PV) and persistent volume claim (PVC) in an ACK cluster.

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

1. Run the following command to deploy the Qwen1.5-4B-Chat model as an inference service by using rtp-llm:

```bash
kubectl apply -f ./yamls/deploy.yaml
kubectl apply -f ./yamls/service.yaml
```

2. Run the following command to view the details of the inference service:

```bash
$ kubectl get po|grep vllm
---
vllm-7cfb9cd9f4-w2hk8   1/1     Running   0          7m11s
```

The output indicates that the inference service is running as expected and is ready to provide services.

### Step 3: Verify the inference service

1. Run the following command to create a port forwarding rule between the inference service and the local environment:

```bash
kubectl port-forward svc/vllm-service 8000:8000
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
{"id":"cmpl-e3efa23d803349d2ad6b44411811109f","object":"chat.completion","created":1720161589,"model":"qwen","choices":[{"index":0,"message":{"role":"assistant","content":"好的，请问您需要测试什么？"},"logprobs":null,"finish_reason":"stop","stop_reason":null}],"usage":{"prompt_tokens":21,"total_tokens":30,"completion_tokens":9}}
```

### (Optional) Step 4: Delete the environment

If you no longer need the resources, delete the resources at the earliest opportunity.

- Run the following command to delete the inference service:

```bash
kubectl delete -f ./yamls
```

## Release Tag

|  tag  |  Date   | release |
|:-----:|:-------:|:-------:|
| 0.4.1 | 2024-07 |  init   |           


