user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;
    sendfile        on;
    keepalive_timeout  65;
    gzip  on;

    include /etc/nginx/conf.d/*.conf;

    (template "upstreams" .)

    (range .)
    server  {
        listen (.FrontendPort);
        server_name (.Domain);
        location / {
            proxy_pass http://(upstreamName .App .Service)(.BackendRootPath);
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }
    (if .EnableSsl)
    server  {
        listen (.SslPort) ssl;
        server_name (.Domain);
        ssl_certificate     (.SslCertPath);
        ssl_certificate_key (.SslKeyPath);
        ssl_session_cache   shared:SSL:10m;
        ssl_session_timeout 10m;
        
        location / {
            proxy_pass http://(upstreamName .App .Service)(.BackendRootPath);
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }
    (end)
    (end)
}
