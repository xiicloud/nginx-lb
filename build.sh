#!/bin/bash

set -e

export CON_NAME=nginx_lb_t
export REG_URL=d.nicescale.com:5000
export IMAGE="nginx-lb"
export TAGS="1.8 1.8.0"
export BASE_IMAGE=microimages/nginx

docker pull $BASE_IMAGE

docker build -t microimages/$IMAGE .

#./test.sh

echo "---> Starting push microimages/$IMAGE:$VERSION"

for t in $TAGS; do
  docker tag -f microimages/$IMAGE microimages/$IMAGE:$t
done

docker push microimages/$IMAGE
