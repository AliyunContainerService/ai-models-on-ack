FROM docker.io/jupyter/base-notebook:python-3.11
RUN conda config --add channels http://mirrors.aliyun.com/pypi/simple/ && \
    sed -i "20i c.NotebookApp.token='' " /etc/jupyter/jupyter_server_config.py
CMD ["sh","-c","jupyter notebook --port 8888 --ip=0.0.0.0"]
