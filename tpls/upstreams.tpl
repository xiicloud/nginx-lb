(range .)
($key := printf "%s-%s" .App .Service)
($etcdKey := printf "/%s/ips" $key)
upstream ($key) {
    {{$ips := ls "($etcdKey)"}}
    {{if len $ips | eq 0}}
    server 1.1.1.1:80; # placeholder
    {{else}}
    {{range $k := ls "($etcdKey)"}}
    server {{base $k}}:(.BackendPort);
    {{end}}
    {{end}}
}
(end)