set CGO_ENABLED=1
set GOARCH=386
go build -o dist/gateway.exe  cmd/main.go
cd dist
gateway.exe -c  config.yml