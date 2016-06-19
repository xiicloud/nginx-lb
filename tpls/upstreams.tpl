(define "upstreams")

(range $service := .)
upstream (upstreamName .App $service.Service) {
(range $app := $service.App)
($etcdKey := getLbKey $app $service.Service)
    {{range $ipKey := ls "($etcdKey)"}}
    server {{base $ipKey}}:($service.BackendPort)(ngxBackup $app);
    {{end}}
(end)
}
(end)

(end)