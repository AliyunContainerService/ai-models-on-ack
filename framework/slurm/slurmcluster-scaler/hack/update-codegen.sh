#!/usr/bin/env bash
#表示有报错即退出 跟set -e含义一样
set -o errexit
#执行脚本的时候，如果遇到不存在的变量，Bash 默认忽略它 ,跟 set -u含义一样
set -o nounset
# 只要一个子命令失败，整个管道命令就失败，脚本就会终止执行 
set -o pipefail
  
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

THIS_PKG="github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler"
source "${CODEGEN_PKG}/kube_codegen.sh"
kube::codegen::gen_helpers \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    "${SCRIPT_ROOT}/apis"
kube::codegen::gen_client \
    --with-watch \
    --output-dir "${SCRIPT_ROOT}/generated" \
    --output-pkg "${THIS_PKG}/generated" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    "${SCRIPT_ROOT}/apis"