#!/usr/bin/env bash
tagName='ccr.ccs.tencentyun.com/g-docker/envoy'

echo "制作镜像..."${tagName}
docker build -t  ${tagName} -f envoy.Dockerfile --no-dec-cache .

if [[ $? != 0 ]]; then
  echo "制作镜像失败"
  exit
fi

if [[ $1 == "push" ]];then
    echo "将镜像推送到云端.."
    docker push  ${tagName}
fi

