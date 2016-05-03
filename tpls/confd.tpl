[template]
prefix      = "/lb/backends"
keys        = (.)
owner       = "root"
mode        = "0644"
src         = "nginx.tpl"
dest        = "/etc/nginx/nginx.conf"
check_cmd   = "/usr/sbin/nginx -t -c {{.src}}"
reload_cmd  = "/usr/sbin/nginx -s reload"
