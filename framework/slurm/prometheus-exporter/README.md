# Slurm Exporter

1. Build exporter 
```bash
FROM golang:1.22 as exporterBuilder
ADD https://api.github.com/repos/KunWuLuan/prometheus-slurm-exporter/git/refs/heads/master version.json
RUN git clone https://github.com/KunWuLuan/prometheus-slurm-exporter.git
RUN cd prometheus-slurm-exporter && go build
```

2. Copy exporter to image

3. Start the exporter when slurmctld start
