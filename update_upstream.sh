#!/bin/bash

if [ -z "$DOMAINS" ]; then
  echo "domains empty, do nothing!"
else
  echo "domains: $DOMAINS"
  for domain in $DOMAINS; do
    config_path=/etc/nginx/conf.d/${domain}.conf

    {
    while [ true ]; do
      backends=$(dig +short $domain|sort)
      backends_config=
      for b in $backends; do
        backends_config=$(echo -e "${backends_config}\nserver $b;")
      done
      cat << EOF > ${config_path}.tmp
upstream backends  {
  $backends_config
}

server {
  listen       80;
  server_name  ${domain};

  location / {
    proxy_pass  http://backends;
  }

  error_page   500 502 503 504  /50x.html;
  location = /50x.html {
    root   /usr/share/nginx/html;
  }
}
EOF
      if [ -f "$config_path" ]; then
        diff $config_path ${config_path}.tmp && sleep 1 && continue
        echo file changed
        mv ${config_path}.tmp $config_path 
        nginx -s reload
      else
        mv ${config_path}.tmp $config_path
        nginx -s reload
      fi

      sleep 2

    done
    } &

  done
fi

exec nginx
