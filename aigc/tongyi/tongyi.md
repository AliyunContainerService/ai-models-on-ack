# 通义千问

## Introduction
https://help.aliyun.com/zh/dashscope/developer-reference/api-details

### Prerequisites
1. [开通DashScope并创建API-KEY](https://help.aliyun.com/zh/dashscope/developer-reference/activate-dashscope-and-create-an-api-key?spm=a2c4g.11186623.0.i2)
2. 申请通义千问模型API权限，[点此申请](https://help.aliyun.com/zh/dashscope/support/faq?spm=a2c4g.11186623.0.i23#vuoFh)

### Deploy Jupyter
```bash
# 1. 创建一个Jupyter Notebook
kubectl apply -f notebook.yaml

# 2. wait deployment ready

# 3. get ExternalIP
kubectl get svc notebook-svc

# 4. open http://${ExternalIP}:8888 

# 5. 新建一个Notebook，参考tongyi.ipynb调用通义千问API
```