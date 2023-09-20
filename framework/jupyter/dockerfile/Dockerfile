FROM yunqi-registry.cn-shanghai.cr.aliyuncs.com/lab/cuda:12.1.0-cudnn8-runtime-ubuntu22.04-py310

RUN pip3 install jupyter && mkdir /root/.jupyter && mkdir /root/notebooks
ADD jupyter_notebook_config.py /root/.jupyter/jupyter_notebook_config.py

WORKDIR /root/notebooks

CMD ["sh","-c","jupyter lab --port 8888 --ip=0.0.0.0 --allow-root"]
