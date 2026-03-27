{{- define "portfolio-gateway.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "portfolio-gateway.fullname" -}}
{{- printf "%s" (include "portfolio-gateway.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "portfolio-gateway.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "portfolio-gateway.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "portfolio-gateway.selectorLabels" -}}
app.kubernetes.io/name: {{ include "portfolio-gateway.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
