FROM microimages/nginx

MAINTAINER william <william@nicescale.com>

RUN mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.bak

ADD update_upstream.sh /update_upstream.sh

CMD ["update_upstream.sh"]
