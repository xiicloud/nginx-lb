(define "upstreams" -)

(range $server := . -)
(range $path, $route := $server.Routes)
# app: {{ $route.App }}, service: {{ $route.Service }}, backend path: {{ $route.BackendPath }}
upstream (upstreamName $route $server.DomainName $path) {
  ($route.UpstreamOptions)
  ($services := split "," $route.Service -)
  (- range $srv := $services -)
  {{range $ipKey := ls "(getLbKey $route $srv)"}}
    server {{base $ipKey}}:($route.Port);
  {{end}}
  (- end)
  (- if $route.Backup -)
  (- $services := split "," $route.Service -)
  (- range $srv := $services -)
  {{range $ipKey := ls "(getLbKey $route.Backup $srv)"}}
    server {{base $ipKey}}:($route.Backup.Port) backup;
  {{end}}
  (- end -)
  (- end -)
}
(end)
(- end )
(- end)
