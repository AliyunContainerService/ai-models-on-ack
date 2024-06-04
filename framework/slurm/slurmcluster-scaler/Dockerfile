FROM golang:1.22 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler
COPY ./ ./

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o slurmctldcluster-scaler main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.cn-hangzhou.aliyuncs.com/acs/alpine:3.18-update
WORKDIR /
COPY --from=builder /usr/local/go/src/github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/slurmctldcluster-scaler .
#USER 65532:65532

ENTRYPOINT ["/slurmctldcluster-scaler"]