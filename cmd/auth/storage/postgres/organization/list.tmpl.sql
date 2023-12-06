SELECT id, legal_name, create_time, update_time
FROM organizations
{{- if .Predicates }}
WHERE {{- range $i, $v := .Predicates }}
	{{- if $i}} AND {{- end }} {{$v -}}
{{- end }}
{{- end -}}
{{- if .Sorting }}
ORDER BY {{- range $i, $v := .Sorting }}
		{{- if $i}}, {{- end }} {{$v.Field }} {{$v.Direction -}}
	{{- end }}
{{- end -}}
{{- if .LimitParam }}
LIMIT {{.LimitParam -}}
{{- end -}}
{{- if .OffsetParam }}
OFFSET {{.OffsetParam -}}
{{- end -}};