#!/usr/bin/env bash
srcDir=${GOPATH}"/src/github.com/gw123/GMQ"
entryFile=cmd/main.go
dstDir=${GOPATH}"/bin/GMQ"
dstName=gateway

if [[ ! -d "${dstDir}" ]]; then
    mkdir -p ${dstDir}
fi

term() {
    echo "term"
    ps -aux|grep er|grep -v grep|awk '{print $2}'|xargs kill -TERM
}

gateway() {
    buildStr="go build -o  ${dstDir}/${dstName}  ${srcDir}/${entryFile}"
    echo ${buildStr}
    ${buildStr}
    if [[ $? -eq 0 ]]; then
      echo "编译成功 开始运行"
      cd ${dstDir} && ./${dstName}
    else
       echo "编译失败"
    fi
}

case $1 in
  "term")
    term
  ;;
  "gateway")
    gateway
  ;;
  *)
    echo "请输入要执行的命令 gateway"
  ;;
esac
