export GOARCH=amd64
export GOOS=linux
go build -o server_linux server.go
upx -9 server_linux