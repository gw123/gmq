#!/usr/bin/env bash
runDocker="docker run -it --rm  \
--network=host \
-e COMMENT_ADDR=0.0.0.0:9090 \
-e DB_HOST=host.docker.internal \
-e DB_USER=gw \
-e DB_PWD=gao123456 \
ccr.ccs.tencentyun.com/g-docker/gateway"

echo ${runDocker}
${runDocker}
