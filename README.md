# 工作原理
基于nginx v1.11.1镜像和confd实现，此LB必须运行在容器里。

1. 通过容器的入口进程自动生成`/etc/confd/conf.d/nginx.toml`，该程序读取`/etc/csphere-services.json`获取需要配置哪些`upstream`和`server_name`
2. 通过配置中心管理`/etc/csphere-services.json`
3. `/etc/csphere-services.json`更新时，执行`kill -USR2 1`重启`confd`进程

# /etc/csphere-services.json格式

```json
{
  "servers": [
    {
      "domain_name": "demo.example.com",
      "frontend_port": 80,
      "opaque": "access_log off;",
      "routes": {
        "/": {
          "app": "myapp",
          "backend_path": "/v2",
          "opaque": "if ( $uri ~ /404 ) {\nreturn 404;\n}",
          "backup": {
            "app": "test",
            "service": "api"
          },
          "port": 80,
          "service": "api"
        },
        "/suburl/": {
          "app": "test2",
          "service": "api"
        }
      },
      "ssl_certificate": "-----BEGIN CERTIFICATE-----\nMIICKTCCAZICCQCBO2ekFdKngDANBgkqhkiG9w0BAQsFADBZMQswCQYDVQQGEwJD\nTjEQMA4GA1UECAwHQmVpamluZzEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQ\ndHkgTHRkMRUwEwYDVQQDDAxwbWEudGVzdC5jb20wHhcNMTYwNTAzMDQxOTUyWhcN\nMTcwNTA0MDQxOTUyWjBZMQswCQYDVQQGEwJDTjEQMA4GA1UECAwHQmVpamluZzEh\nMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMRUwEwYDVQQDDAxwbWEu\ndGVzdC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALo7afOGNwyXqN24\nViyG+HQG4In0O4SqrAQr6WUOg2g95rf7qSEG6U4o+wO5BM8nd6wzW0JnNqTxKXpf\nnV2Ebub0uoITUHtbFMRSEYfLShQqbKBltnk09P9p4IcVhM5vUd+G9reGagaH84bt\nP2bSZ5JmKBsiUf277b4ZlM/nIS4rAgMBAAEwDQYJKoZIhvcNAQELBQADgYEAOp9q\noNF/sGwEJzUKkrZ9jNfr9nvcVwsR9VajsB0cQW859IQ75r0P80NwPwJ4qbIMNsid\n/1qwqzZlnYvOm01176DQTCgRC42r4vbLMKzNR4Xf6O25gwWaa2ZSSpLNQLxauKSg\n2h/qgz7Rfn7rYMYZmVLNnnJqujr8GbZ4cKyuy1w=\n-----END CERTIFICATE-----",
      "ssl_certificate_key": "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQC6O2nzhjcMl6jduFYshvh0BuCJ9DuEqqwEK+llDoNoPea3+6kh\nBulOKPsDuQTPJ3esM1tCZzak8Sl6X51dhG7m9LqCE1B7WxTEUhGHy0oUKmygZbZ5\nNPT/aeCHFYTOb1Hfhva3hmoGh/OG7T9m0meSZigbIlH9u+2+GZTP5yEuKwIDAQAB\nAoGAbNZeRFkzAOP9Z57cledHep+uSFF5Gz6Xi1SScWH7AEf0959XJ5sfbHNcx78w\nhVR+hx/4fKVPdTQP1pncoRPNr6GcK6+9NhURiBy4oaWIAlswvuMeCFBJp4KBI3np\n/hQIXKHZ4hNasD5SzHBo5bJOG6P5577KD4t39QFcBMGy8xECQQDz77IbLXOG+UWJ\nSpw84HUJxzWAc7h9IgjkybJQIDoHvgmsYIjLFnLimwIbyvMlPBnyAopvZoO5tabX\nWfGmQHJpAkEAw3Eqe7XG2j023SUDe8slqOIlnivdP95M2tACOKIvbJFB83zHDK2C\n46cdyoqfI/bLimnxgPxYMq5CTr33I7RhcwJBAJTlU17ZcHILx5EU1KcoDuiICzU7\n7XmcA7e7Ebds5F8DdZ4dUoI8UqXVHgVe7OlmdSPOvzdeaLs7kPpUMXdcUTkCQGG5\numZ1dGM37LETivRhlgkmW20FvfHrtD5NeG7dGh2NXI7lu5opQKOYsprOSdjv1ML3\nSp0WkPt2iw1Yi7U8wuUCQQCh3Z/svnkDSxBrexFUwt5RLhF5YBJUtIgeOBQ/wlDp\nhQPBEwDsQoM9LxnEmVNyzeF8Yz9RJ6gANpnZYLszWrGk\n-----END RSA PRIVATE KEY-----",
      "ssl_port": 443
    }
  ],
  "version": "0.1"
}
```

`servers`这个数组里各字段的说明：

- `domain_name` 对外提供服务的域名
- `frontend_port` 负载均衡前端监听的端口，默认为80
- `opaque` 用户自定义的需要放到nginx 配置文件`server`里面的配置项
- `routes` 转发规则，key表示供前端访问的URL前缀，value表示对应的后端
- `ssl_certificate` SSL证书的内容
- `ssl_certificate_key` SSL证书密钥
- `ssl_port` SSL监听端口

`routes`下面每一条记录(route)的字段说明：

- `app` 应用名字，在应用详情页面能查看到
- `backend_path` 后端服务根目录的路径，默认为`/`
- `backup` 如果配置了这一项，会把这里面配置的服务当作nginx upstream配置里的backup servers. 其结构与routes下面的每一条记录结构一致
- `opaque` 用户自定义的需要放到nginx 配置文件`location`里面的配置项
- `port` 后端服务的端口，默认为80
- `service` 应用里的服务的名字，在应用详情页面能查看到


如果etcd中设置以下key: **(如果是使用cSphere平台，创建相关应用即可，无需手工设置以下key)**

```
etcdctl set /lb/backends/myapp-api/ips/192.168.122.57 '{"weight":50}'
etcdctl set /lb/backends/test-api/ips/192.168.122.58 '{"weight":50}'
etcdctl set /lb/backends/test2-api/ips/192.168.122.59 '{"weight":50}'
```

**weight参数目前暂时没有用到**

通过以上示例配置文件，将生成如下的nginx配置文件：

```
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

    
upstream myapp-api-test-api {
  
    server 192.168.122.57:80;
  
    server 192.168.122.58:80 backup;
  }

upstream test2-api {
  
    server 192.168.122.59:80;
  }

    server  {
        listen 80;
        server_name demo.example.com;
        access_log off;

        
        
        location / {
            if ( $uri ~ /404 ) {
return 404;
}
            proxy_pass http://myapp-api-test-api/v2/;
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
        
        location /suburl/ {
            
            proxy_pass http://test2-api/;
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }
    
    server  {
        listen 443 ssl;
        server_name demo.example.com;
        ssl_certificate     /etc/nginx/ssl/demo.example.com.pem;
        ssl_certificate_key /etc/nginx/ssl/demo.example.com.key;
        ssl_session_cache   shared:SSL:10m;
        ssl_session_timeout 10m;
        access_log off;
        
        location / {
            if ( $uri ~ /404 ) {
return 404;
}
            proxy_pass http://myapp-api-test-api/v2/;
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
        
        location /suburl/ {
            
            proxy_pass http://test2-api/;
            proxy_redirect    off;
            proxy_set_header  Host             $host;
            proxy_set_header  X-Real-IP        $remote_addr;
            proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }
    
}
```
