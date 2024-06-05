#!/bin/bash
OS=windows
if [ $1 ]
then
    OS=$1
fi
BIN_DIR=../../bin
function build() {
    PROJECT=$1
    if [ $OS = "windows" ]
    then
        if [ -f "$BIN_DIR/$PROJECT.exe" ]
        then
            return
        fi
    elif [ -f "$BIN_DIR/$PROJECT" ]
        then
        return
    fi
    VERSION_PATH=github.com/dxasu/tools/lancet/version
    VER=$(echo "No.$(git log --oneline |wc -l)" | sed 's/ //g')
    TAG=$(git log -n1 --pretty=format:%h |git tag --contains)
    # STATUS="-X '$VERSION_PATH.GitStatus=$(git status)'"
    LDFLAG="-s -X '$VERSION_PATH.Version=$VER $TAG'  -X '$VERSION_PATH.GitCommit=$(git rev-parse --short HEAD)' -X '$VERSION_PATH.BuildTime=$(date "+%Y-%m-%d %H:%M:%S")'"
    CGO_ENABLED=0 GOOS=$OS GOARCH=amd64 go build -ldflags "$LDFLAG" -o $BIN_DIR
    echo $PROJECT
# windows 下编译
# SET CGO_ENABLED=0
# SET GOOS=darwin
# SET GOARCH=amd64
# go build
}

cd ../cmd
for dir in `ls`
do
    if [ -d $dir ] 
    then
        cd $dir
        build $dir
        cd ..
    fi
done
