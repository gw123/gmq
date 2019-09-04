#!/usr/bin/env bash
#/root/code/output/mfrpc -c /root/code/output/frpc.ini
outputDir=/root/code/output/
GMQ=gateway

export GO111MODULE=on


function build() {
    echo "build ${GMQ}..."
    go build -o ${outputDir}${GMQ} cmd/main.go
    if [ $? == 0 ]
    then
        echo "build ${GMQ} over... ,dist"${outputDir}${GMQ}
    else
        echo "build faild"
    fi
}


case $1 in
   build)
   build
   ;;
esac
