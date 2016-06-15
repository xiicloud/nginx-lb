# 工作原理
1. 通过容器的入口进程自动生成`/etc/confd/conf.d/nginx.toml`，该程序读取`/etc/csphere-services.json`获取需要配置哪些`upstream`和`server_name`
2. 通过配置中心管理`/etc/csphere-services.json`
3. `/etc/csphere-services.json`更新时，执行`kill -USR2 1`重启`confd`进程

# /etc/csphere-services.json格式

```json
[
  {
    "domain_name": "demo.example.com",
    "app": "demo",
    "service": "app",
    "backend_port": 8080,
    "frontend_port": 80,
    "backend_root_path": "/",
    "ssl_certificate": "SSL certificate content",
    "ssl_certificate_key": "SSL certificate key content",
    "ssl_port": 443
  },
  {
    "domain_name": "test.example.com",
    "app": ["test", "test2"],
    "service": "app",
    "backend_port": 8080,
    "frontend_port": 80,
    "backend_root_path": "/myapp"
  }
]
```

- `domain_name` 对外提供服务的域名
- `app` 应用名字，在应用详情页面能查看到，可以支持多个应用名，用来做粗暴的灰度发布
- `service` 应用里的服务的名字，在应用详情页面能查看到
- `backend_port` 后端服务的端口，默认为80
- `frontend_port` 负载均衡前端监听的端口，默认为80

以上配置文件示例表明：

- 当用户访问`demo.example.com`时，代理到`demo`应用中的`app`服务
- 当用户访问`test.example.com`时，代理到`test`应用中的`app`服务

# etcd中的数据格式
```
etcdctl set /lb/backends/demo-app/ips/192.168.58.57 '{"weight":50}'
etcdctl set /lb/backends/ci-gitlab/ips/192.168.57.74 '{"weight":50}'
etcdctl set /lb/backends/ci-archiva/ips/192.168.57.55 '{"weight":50}'
```

**weight参数目前暂时没有用到**

# 支持的环境变量
支持通过`MODE`环境变量来决定运行模式。内置了三种运行模式：

- `default` 全自动配置，自动生成upstream和server配置文件，每个upstream对应一个独立的域名
- `path-mux` 全自动配置，所有服务通过IP访问，不同服务通过URL的前缀进行区分
