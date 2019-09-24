#!/usr/bin/env bash
dockerImage=golang:1.12.9-alpine3.10
srcDist=entry/server.go
export dstExe=dist/gateway
tagName='ccr.ccs.tencentyun.com/g-docker/gateway'

runStr="docker run -it  --rm -v $GOPATH:/go \
-v $HOME/.ssh:/root/.ssh \
-v $PWD:/usr/src/myapp \
-w /usr/src/myapp \
$dockerImage \
go build -v -o $dstExe /usr/src/myapp/main.go"


echo "编译代码........."
echo ${runStr}
${runStr}

if [[ $? != 0 ]]; then
  echo "编译失败"
  exit
fi

echo "制作镜像........."
docker build -t  ${tagName}  --no-cache .

if [[ $? != 0 ]]; then
  echo "制作镜像失败"
  exit
fi

if [[ $1 == "push" ]];then
    echo "将镜像推送到云端.."
    docker push  ${tagName}
fi

