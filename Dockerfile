FROM microimages/nginx

MAINTAINER william <william@nicescale.com>

RUN mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.bak && \
  mkdir /etc/service.d/nginx/ && \
  echo 'nginx -g "daemon off;"' > /etc/service.d/nginx/run && \
  chmod 755 /etc/service.d/nginx/run 

ADD update_upstream.sh /update_upstream.sh

CMD ["/update_stream.sh"]
