FROM nginx
MAINTAINER Shijiang Wei<mountkin@gmail.com>

RUN apt-get update && \
  apt-get install -y --no-install-recommends curl && \
  rm -fr /var/lib/apt/lists/*

RUN curl -L -o /bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64 && \
  chmod +x /bin/confd

RUN rm /etc/nginx/conf.d/default.conf
ADD init/init /init
ADD tpls /tpls
ADD csphere-services.json /etc/

ENTRYPOINT ["/init"]
EXPOSE 80 443
