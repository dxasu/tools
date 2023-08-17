# tools
collection effective tools

# go install 失败，代理设置
go env -w GOPROXY=https://goproxy.cn

# make
make -f ../../scripts/Makefile

## qrcode
go install -ldflags="-s" github.com/dxasu/tools/cmd/qrcode@latest

## clipboard
go install -ldflags="-s" github.com/dxasu/tools/cmd/clipboard@latest

## otp
go install -ldflags="-s" github.com/dxasu/tools/cmd/otp@latest