# tools
collection effective tools

# go install 失败，代理设置
go env -w GOPROXY=https://goproxy.cn

# --ldflags
--ldflags="-s -X 'github.com/dxasu/tools/lancet/version.Version=$(echo "No.$(git log --oneline |wc -l)" | sed 's/ //g')$(git log -n1 --pretty=format:%h |git tag --contains)'  -X 'github.com/dxasu/tools/lancet/version.GitCommit=$(git rev-parse --short HEAD)' -X 'github.com/dxasu/tools/lancet/version.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'github.com/dxasu/tools/lancet/version.GitStatus=$(git status)'"


# make
cd cmd/qrcode
make -f ../../scripts/Makefile

## qrcode
go install -ldflags="-s" github.com/dxasu/tools/cmd/qrcode@latest

## clipboard
go install -ldflags="-s" github.com/dxasu/tools/cmd/clipboard@latest

## jsonhand
go install -ldflags="-s" github.com/dxasu/tools/cmd/jsonhand@latest

## otp
go install -ldflags="-s" github.com/dxasu/tools/cmd/otp@latest

## rsa
#私钥
openssl genrsa 2048 | openssl pkcs8 -topk8 -nocrypt -out private.key.pem
#公钥
openssl rsa -in private.key.pem -pubout > public.key.pem