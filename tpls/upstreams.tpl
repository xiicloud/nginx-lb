(define "upstreams" -)

(range $server := . -)
(range $route := $server.Routes)
upstream (upstreamName $route) {
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
