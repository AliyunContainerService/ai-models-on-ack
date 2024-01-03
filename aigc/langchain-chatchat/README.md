# Langchain-chatchat

基于开源本地知识库问答项目[Langchain-Chatchat](https://github.com/chatchat-space/Langchain-Chatchat),
集成了阿里云推理引擎DeepGPU-LLM，AnalyticDB for PostgreSQL向量数仓等产品，快速构建检索增强生成(RAG)大模型知识库项目。

## 注意事项

<font color="red">
阿里云不对第三方模型的合法性、安全性、准确性进行任何保证，阿里云不对由此引发的任何损害承担责任。

您应自觉遵守第三方模型的用户协议、使用规范和相关法律法规，并就使用第三方模型的合法性、合规性自行承担相关责任。
</font></br>

## 前提条件

- 已创建ACK
  GPU集群。具体操作，请参见[如何创建GPU集群](https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/user-guide/create-an-ack-managed-cluster-with-gpu-accelerated-nodes)。
- GPU机型需要为V100、A10系列，不支持vgpu。为避免显存不足，建议GPU机型显存大于等于24G。

## 组件说明

Langchain-chatchat由Chat及LLM组成。

- Chat：ChatBot应用，支持LLM问答和知识库问答。Database用于存储知识库Embedding后的向量。Database支持`faiss`及`ADB`两种类型。
- LLM：模型推理服务，基于开源的FastChat项目部署LLM模型。默认使用DeepGPU-LLM加速推理的qwen-7b-chat-aiacc模型。支持替换为其他开源模型或DeepGPU-LLM加速模型，替换模型可通过pvc挂载到容器中。

## 访问服务

1. 登录容器服务管理控制台，在左侧导航栏选择集群。
2. 在集群列表页面，单击目标集群名称，选择工作负载 > 无状态。
3. 在无状态列表页面，选择Langchain-chatchat部署的命名空间，等待Pod就绪，即容器组数量变为1/1。(
   注意：镜像拉取大约需要10-20分钟)
4. 在左侧导航栏，选择网络 > 服务。查看Langchain-chatchat部署的服务名称，service名称格式为chat-{releaseName}。
5. 将服务转发到本地，将占用本地的8501端口，参考[port-forward](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster)。

```bash
# 将chat-langchain替换为具体service的名称
kubectl port-forward service/chat-langchain 8501:8501
```

命令执行成功后，输出内容如下。

```text
Forwarding from 127.0.0.1:8501 -> 8501
Forwarding from [::1]:8501 -> 8501
```

6. 在本地通过浏览器访问服务[http://localhost:8501](http://localhost:8501)。

## 应用配置

### 参数配置说明

| 参数                    | 描述                                                                          | 默认值                        |
|:----------------------|:----------------------------------------------------------------------------|:---------------------------|
| llm.model             | llm模型名称                                                                     | qwen-7b-chat-aiacc         |
| llm.load8bit          | llm模型int8量化                                                                 | true                       |
| llm.modelPVC          | 模型存储PVC，挂载到容器内/llm-model目录                                                  | true                       |
| llm.pod.replicas      | 模型推理服务副本数量                                                                  | 1                          |
| llm.pod.instanceType  | 模型推理服务部署方式，取值：<br/>ecs：部署到ECS节点上。<br/>eci： 部署到ECI上（ACK Serverless集群请使用eci）。 | ecs                        |
| chat.pod.replicas     | 应用服务副本数量                                                                    | 1                          |
| chat.pod.instanceType | 应用部署方式，取值：<br/>ecs：部署到ECS节点上。<br/>eci： 部署到ECI上（ACK Serverless集群请使用eci）。     | ecs                        |
| chat.kbPVC            | 已存在的PVC，用于保存本地知识库文件                                                         | 无                          |
| db.dbType             | 向量数据库类型，支持faiss, adb                                                        | faiss                      |
| db.embeddingModel     | embedding模型                                                                 | text2vec-bge-large-chinese |

### 服务端口说明

| 组件                   | 端口         | 描述              |
|:---------------------|:-----------|:----------------|
| chat                 | 8501       | Web UI 页面       |
| chat                 | 7861       | FastAPI Docs 页面 |
| llm                  | 8888       | 模型推理接口          |

### 常见修改配置操作

#### 采用ECI Pod方式部署

修改参数配置，将`llm.pod.instanceType`和`chat.pod.instanceType`设置为`eci`。  
ECI类型默认的Annotation及Label配置如下。您可以参考文档[ECI Pod Annotation](https://help.aliyun.com/zh/eci/user-guide/pod-annotations-1)
查看Annotation的详细信息。

```yaml
...
annotations:
  k8s.aliyun.com/eci-use-specs: ecs.gn6i-c8g1.2xlarge,ecs.gn6i-c16g1.4xlarge,ecs.gn6v-c8g1.8xlarge,ecs.gn7i-c8g1.2xlarge,ecs.gn7i-c16g1.4xlarge
  k8s.aliyun.com/eci-extra-ephemeral-storage: "50Gi"
labels:
  alibabacloud.com/eci: "true"
...
```

如果您更改了镜像或者模型，则需要修改`k8s.aliyun.com/eci-use-specs`及`k8s.aliyun.com/eci-extra-ephemeral-storage`
annotation, 否则会导致应用因为资源不足无法启动。

> ECI Pod方式部署不支持DeepGPU-LLM，需使用开源模型。

#### 修改副本数量

设置`llm.pod.replicas`为需要的推理服务副本数量。  
设置`chat.pod.replicas`为需要的应用服务副本数量。

### 向量数据库

#### faiss

faiss是由facebook开源的一款内存向量库，项目地址[https://github.com/facebookresearch/faiss](https://github.com/facebookresearch/faiss)。  
faiss内存数据库部署在chat pod中，受chat pod的资源约束。如果使用faiss向量数据库，建议增加chat pod的内存资源。

#### AnalyticDB PostgreSQL (简称 ADB)

云原生数据仓库AnalyticDB
PostgreSQL版是一种大规模并行处理（MPP）数据仓库服务，可提供海量数据在线分析服务。产品简介请参考文档[什么是云原生数据仓库](https://help.aliyun.com/zh/analyticdb-for-postgresql/product-overview/overview-product-overview)。  
Langchain-chatchat项目中使用的ADB需要满足以下条件：

- 需开启向量引擎优化功能
- 计算节点规格>=4C16G

### embedding模型

应用内置的embedding模型为text2vec-bge-large-chinese，详情请参考[hugging face文档](https://huggingface.co/shibing624/text2vec-bge-large-chinese)。

chat应用默认使用CPU运行embedding模型，可通过在`chat.pod.resources`中申请GPU资源来提高文本向量化速度。

## 模型配置

### 支持的模型列表

| 模型类型            | 模型名称                     | 容器内模型文件路径                           |
|:----------------|:-------------------------|:------------------------------------|
| DeepGPU-LLM转换模型 | qwen-7b-chat-aiacc       | /llm-model/qwen-7b-chat-aiacc       |
| DeepGPU-LLM转换模型 | qwen-14b-chat-aiacc      | /llm-model/qwen-14b-chat-aiacc      |
| DeepGPU-LLM转换模型 | chatglm2-6b-aiacc        | /llm-model/chatglm2-6b-aiacc        |
| DeepGPU-LLM转换模型 | baichuan2-7b-chat-aiacc  | /llm-model/baichuan2-7b-chat-aiacc  |
| DeepGPU-LLM转换模型 | baichuan2-13b-chat-aiacc | /llm-model/baichuan2-13b-chat-aiacc |
| DeepGPU-LLM转换模型 | llama-2-7b-hf-aiacc      | /llm-model/llama-2-7b-hf-aiacc      |
| DeepGPU-LLM转换模型 | llama-2-13b-hf-aiacc     | /llm-model/llama-2-13b-hf-aiacc     |
| 开源模型            | qwen-7b-chat             | /llm-model/Qwen-7B-Chat             |
| 开源模型            | qwen-14b-chat            | /llm-model/Qwen-14B-Chat            |
| 开源模型            | chatglm2-6b              | /llm-model/chatglm2-6b              |
| 开源模型            | chatglm2-6b-32k          | /llm-model/chatglm2-6b-32k          |
| 开源模型            | baichuan2-7b-chat        | /llm-model/Baichuan2-7B-Chat        |
| 开源模型            | baichuan2-13b-chat       | /llm-model/Baichuan2-13B-Chat       |
| 开源模型            | llama-2-7b-hf            | /llm-model/Llama-2-7b-hf            |
| 开源模型            | llama-2-13b-hf           | /llm-model/Llama-2-13b-hf           |


### 使用DeepGPU-LLM转换模型

[DeepGPU-LLM](https://help.aliyun.com/zh/egs/what-is-deepgpu-llm)是阿里云研发的基于GPU云服务器的大语言模型（Large Language Model，LLM）推理引擎，旨在优化大语言模型在GPU云服务器上的推理过程，通过优化和并行计算等技术手段，提供免费的高性能、低延迟推理服务。DeepGPU-LLM使用方式请参考文档[使用DeepGPU-LLM实现大语言模型在GPU上的推理优化](https://help.aliyun.com/zh/egs/developer-reference/install-and-use-deepgpu-llm)。  

Langchain-chatchat项目已安装DeepGPU-LLM，默认使用DeepGPU-LLM加速后的模型qwen-7b-chat-aiacc。 

如想要使用DeepGPU-LLM对其他开源LLM模型进行推理优化，您需要先将huggingface格式的开源模型转换为DeepGPU-LLM支持的格式，然后才能使用DeepGPU_LLM进行模型的推理优化服务。以qwen-7b-chat为例，可使用如下命令在容器中进行模型格式转换：  
```text
#qwen-7b weight convert
huggingface_qwen_convert \
    -in_file /llm-model/Qwen-7B-Chat \
    -saved_dir /llm-model/qwen-7b-chat-aiacc \
    -infer_gpu_num 1 \
    -weight_data_type fp16 \
    -model_name qwen-7b
```

### 更换模型
#### 步骤一：创建静态PV及PVC
- OSS模型
1. 执行以下命令创建Secret。
```bash
kubectl create -f oss-secret.yaml
```
以下为创建Secret的oss-secret.yaml示例文件，需要指定akId和akSecret。
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: oss-secret
  namespace: default
stringData:
  akId: <your AccessKeyID>
  akSecret: <your AccessKeySecret>
```
2. 执行以下命令创建静态卷PV。
```bash
  kubectl create -f model-oss.yaml
```
以下为创建静态卷PV的model-oss.yaml示例文件，需要指定bucket,url等参数。
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: model-oss
  labels:
    alicloud-pvname: model-oss
spec:
  capacity:
    storage: 30Gi 
  accessModes:
    - ReadOnlyMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: ossplugin.csi.alibabacloud.com
    volumeHandle: model-oss
    nodePublishSecretRef:
      name: oss-secret
      namespace: default
    volumeAttributes:
      bucket: "<your bucket name>"
      url: "<your oss endpoint>" # oss-cn-hangzhou.aliyuncs.com
      otherOpts: "-o umask=022 -o max_stat_cache_size=0 -o allow_other"
      path: "/"
```
3. 执行以下命令创建静态卷PVC。
```bash
kubectl create -f pvc-oss.yaml
```
以下为创建静态卷PVC的model-pvc.yaml示例文件。
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: model-pvc
spec:
  accessModes:
    - ReadOnlyMany
  resources:
    requests:
      storage: 30Gi
  selector:
    matchLabels:
      alicloud-pvname: model-oss
```
参数配置可参考[使用OSS静态存储卷](https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/user-guide/mount-statically-provisioned-oss-volumes)。
- NAS模型
1. 执行以下命令创建静态卷PV。
```bash
kubectl create -f model-nas.yaml
```
以下为创建静态卷PV的model-nas.yaml示例文件，需要指定NAS服务地址和路径。
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: model-nas
  labels:
    alicloud-pvname: model-nas
spec:
  capacity:
    storage: 30Gi
  accessModes:
    - ReadWriteMany
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: model-nas 
    volumeAttributes:
      server: "<your nas server>"
      path: "<your model path>"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3
```
2. 执行以下命令创建静态卷PVC。
```bash
kubectl create -f model-pvc.yaml
```
以下为创建静态卷PVC的model-pvc.yaml示例文件。
```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: model-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 30Gi
  selector:
    matchLabels:
      alicloud-pvname: model-nas
```
参数配置可参考[使用NAS静态存储卷](https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/user-guide/mount-statically-provisioned-nas-volumes)。

#### 步骤二：更新Helm Value
1. 登陆容器服务控制台，选择对应的ACK集群。
2. 点击左侧菜单栏，选择应用->Helm，找到部署的LLM服务
3. 点击右侧更新按钮，更改`llm.model`为新的模型名称，`llm.modelPVC`为存储新模型的pvc名称。模型名称及模型挂载路径参考支持的模型列表。


## Release Note

| 版本号     | 变更时间         | 变更内容                                               |
|---------|--------------|----------------------------------------------------|
| `0.1.0` | 2023年12月26日	 | 支持阿里云推理引擎DeepGPU-LLM，AnalyticDB for PostgreSQL向量数仓 |