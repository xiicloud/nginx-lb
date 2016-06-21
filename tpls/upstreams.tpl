(define "upstreams" -)

(range $server := . -)
(range $route := $server.Routes)
upstream (upstreamName $route) {
  {{range $ipKey := ls "(getLbKey $route)"}}
    server {{base $ipKey}}:($route.Port);
  {{end}}
  (- if $route.Backup -)
  {{range $ipKey := ls "(getLbKey $route.Backup)"}}
    server {{base $ipKey}}:($route.Backup.Port) backup;
  {{end}}
  (- end -)
}
(end)
(- end )
(- end)