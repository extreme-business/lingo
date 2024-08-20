SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time
FROM users u 
{{- if .Predicates }}
WHERE {{- range $i, $v := .Predicates }}
	{{- if $i}} AND {{- end }} {{$v -}}
{{- end }}
{{- end -}}
{{- if .Sorting }}
ORDER BY {{- range $i, $v := .Sorting }}
		{{- if $i}}, {{- end }} u.{{$v.Field }} {{$v.Direction -}}
	{{- end }}
{{- end -}}
{{- if .LimitParam }}
LIMIT {{.LimitParam -}}
{{- end -}}
{{- if .OffsetParam }}
OFFSET {{.OffsetParam -}}
{{- end -}};