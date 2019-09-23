#!/usr/bin/env bash
protoc -I=.  --go_out=plugins=grpc:. comment.proto

protoc -I=.  comment.proto \
 --js_out=import_style=commonjs:. \
 --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.
