# tools
collection effective tools

# install use cmd
支持Windows/Linux/MacOS。安装很简单：

Linux和Mac在终端运行：
```
sh -c "$(curl -fsSL https://raw.githubusercontent.com/dxasu/tools/master/scripts/install.sh)"
```
如果http://raw.githubusercontent.com被屏蔽，改用git clone https://github.com/dxasu/tools && bash kd/install.sh。ArchLinux用户可以直接通过AUR安装（例如`yay -S kd`）。
Win用powershell运行：
```
Invoke-WebRequest -uri 'https://github.com/dxasu/tools/releases/latest/download/kd_windows_amd64.exe' -OutFile ( New-Item -Path "C:\bin\kd.exe" -Force )
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\bin", "User")
```


# go install 失败，代理设置
go env -w GOPROXY=https://goproxy.cn

# --ldflags
--ldflags="-s -X 'github.com/dxasu/pure/version.Version=$(echo "No.$(git log --oneline |wc -l)" | sed 's/ //g')$(git log -n1 --pretty=format:%h |git tag --contains)'  -X 'github.com/dxasu/pure/version.GitCommit=$(git rev-parse --short HEAD)' -X 'github.com/dxasu/pure/version.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'github.com/dxasu/pure/version.GitStatus=$(git status)'"


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