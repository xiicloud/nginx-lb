[template]
prefix      = "/lb/backends"
keys        = (.Keys)
owner       = "root"
mode        = "0644"
src         = "nginx.tpl"
dest        = "(.DestPath)"
check_cmd   = "/usr/sbin/nginx -t -c {{.src}}"
reload_cmd  = "/usr/sbin/nginx -s reload"
