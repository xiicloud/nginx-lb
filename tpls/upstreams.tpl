(define "upstreams")

(range $service := .)
upstream ($service.UpstreamName) {
(range $app := $service.App)
($etcdKey := printf "/%s-%s/ips" $app $service.Service)
    {{range $ipKey := ls "($etcdKey)"}}
    server {{base $ipKey}}:($service.BackendPort);
    {{end}}
(end)
}
(end)

(end)