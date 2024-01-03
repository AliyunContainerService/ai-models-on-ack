{{/*
Returns the chat annotations
*/}}
{{- define "chat.annotations" -}}
  {{- $annotations := dict -}}
  {{- if eq .Values.chat.pod.instanceType "eci" -}}
    {{ $annotations = fromYaml (include "chat.eci.annotations" .) -}}
  {{- end -}}
  {{- if .Values.chat.pod.annotations -}}
    {{ $annotations = mergeOverwrite $annotations .Values.chat.pod.annotations -}}
  {{- end -}}
  {{- if $annotations -}}
    {{- toYaml $annotations | indent 0 -}}
  {{- end -}}
{{- end -}}

{{/*
Returns the llm annotations
*/}}
{{- define "llm.annotations" -}}
  {{- $annotations := dict -}}
  {{- if eq .Values.llm.pod.instanceType "eci" -}}
    {{ $annotations = fromYaml (include "llm.eci.annotations" .) -}}
  {{- end -}}
  {{- if .Values.llm.pod.annotations -}}
    {{ $annotations = mergeOverwrite $annotations .Values.llm.pod.annotations -}}
  {{- end -}}
  {{- if $annotations -}}
    {{- toYaml $annotations | indent 0 -}}
  {{- end -}}
{{- end -}}

{{/*
Returns the chat labels
*/}}
{{- define "chat.labels" -}}
  {{- $labels := dict -}}
  {{- if .Values.chat.pod.labels -}}
    {{ $labels = .Values.chat.pod.labels -}}
  {{- end -}}
  {{- if eq .Values.chat.pod.instanceType "eci" -}}
    {{ $labels = mergeOverwrite $labels (fromYaml (include "eci.labels" .)) -}}
  {{- end -}}
  {{ $labels = mergeOverwrite $labels (fromYaml (include "chat.matchLabels" .)) -}}
  {{- if $labels -}}
    {{- toYaml $labels | indent 0 -}}
  {{- end -}}
{{- end -}}

{{/*
Returns the llm labels
*/}}
{{- define "llm.labels" -}}
  {{- $labels := dict -}}
  {{- if .Values.llm.pod.labels -}}
    {{ $labels = .Values.llm.pod.labels -}}
  {{- end -}}
  {{- if eq .Values.llm.pod.instanceType "eci" -}}
    {{ $labels = mergeOverwrite $labels (fromYaml (include "eci.labels" .)) -}}
  {{- end -}}
  {{ $labels = mergeOverwrite $labels (fromYaml (include "llm.matchLabels" .)) -}}
  {{- if $labels -}}
    {{- toYaml $labels | indent 0 -}}
  {{- end -}}
{{- end -}}

{{/*
Returns if we should generate the nginx.
*/}}
{{- define "chat.enabled.nginx" -}}
    {{- if .Values.nginx.enabled -}}
        {{- true -}}
    {{- end -}}
{{- end -}}


{{/*
Returns the chat match labels
*/}}
{{- define "chat.matchLabels" -}}
app.kubernetes.io/name: {{ .Values.appNameOverride | default (printf "chat-%s" .Chart.Name) | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/instance: chat-{{ .Release.Name }}
{{- end -}}

{{/*
Returns the llm match labels
*/}}
{{- define "llm.matchLabels" -}}
app.kubernetes.io/name: {{ .Values.appNameOverride | default (printf "llm-%s" .Chart.Name) | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/instance: llm-{{ .Release.Name }}
{{- end -}}

{{- define "llm.eci.annotations" -}}
k8s.aliyun.com/eci-use-specs: ecs.gn6i-c8g1.2xlarge,ecs.gn6i-c16g1.4xlarge,ecs.gn6v-c8g1.8xlarge,ecs.gn7i-c8g1.2xlarge,ecs.gn7i-c16g1.4xlarge
k8s.aliyun.com/eci-extra-ephemeral-storage: "100Gi"
{{- end -}}

{{- define "chat.eci.annotations" -}}
k8s.aliyun.com/eci-extra-ephemeral-storage: "100Gi"
{{- end -}}

{{- define "eci.labels" -}}
alibabacloud.com/eci: "true"
{{- end -}}


{{/*
Return the proper image name
*/}}
{{- define "chat.image" -}}
    {{- $registryName := default "kube-ai-registry.cn-shanghai.cr.aliyuncs.com" .Values.chat.pod.image.registry -}}
    {{- $repositoryName := default "kube-ai/chatchat" .Values.chat.pod.image.repository -}}
    {{- $tag := .Values.chat.pod.image.tag | default "aiacc-v0.1.4" | toString -}}
    {{- if not (hasSuffix "-aiacc" .Values.llm.model) }}
    {{- $tag = trimPrefix "aiacc-" $tag -}}
    {{- end -}}
    {{- if hasPrefix "sha256:" $tag }}
    {{- printf "%s/%s@%s" $registryName $repositoryName $tag -}}
    {{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- end -}}

{{/*
Return the proper image name
*/}}
{{- define "llm.image" -}}
    {{- $registryName := default "kube-ai-registry.cn-shanghai.cr.aliyuncs.com" .Values.llm.pod.image.registry -}}
    {{- $repositoryName := default "kube-ai/chatchat" .Values.llm.pod.image.repository -}}
    {{- $tag := .Values.llm.pod.image.tag | default "aiacc-v0.1.4" | toString -}}
    {{- if not (hasSuffix "-aiacc" .Values.llm.model) }}
    {{- $tag = trimPrefix "aiacc-" $tag -}}
    {{- end -}}
    {{- if hasPrefix "sha256:" $tag }}
    {{- printf "%s/%s@%s" $registryName $repositoryName $tag -}}
    {{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- end -}}

{{/*
Return the llm load8bit
*/}}
{{- define "llm.load8bit" -}}
    {{- $load8bit := ne .Values.llm.load8bit false -}}
    {{- if not (hasSuffix "-aiacc" .Values.llm.model) }}
    {{- if hasPrefix "qwen-" .Values.llm.model }}
        {{- $load8bit = "false" -}}
    {{- end -}}
    {{- end -}}
    {{- printf "%v" $load8bit -}}
{{- end -}}
