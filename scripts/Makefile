# 颜色
# RED    := $(shell tput -Txterm setaf 1)
# GREEN  := $(shell tput -Txterm setaf 2)
# YELLOW := $(shell tput -Txterm setaf 3)
# VIOLET := $(shell tput -Txterm setaf 5)
# AQUA   := $(shell tput -Txterm setaf 6)
# WHITE  := $(shell tput -Txterm setaf 7)
# RESET  := $(shell tput -Txterm sgr0)
VERSION_PATH = github.com/dxasu/tools/lancet/version
# 判断当前commit是否有tag如果有tag则显示tag没有则显示 commit次数.哈希
VER = $(shell echo "No.$(shell git log --oneline |wc -l)" | sed 's/ //g')
# 用于判断当前是否有tag
TAG = $(shell git log -n1 --pretty=format:%h |git tag --contains)

VERSION =  -X '$(VERSION_PATH).Version=$(VER) $(TAG)'
# GO_PATH = $(shell cat go.mod |grep module |cut -b 8-)
GIT_COMMIT = -X '$(VERSION_PATH).GitCommit=$$(git rev-parse --short HEAD)'
BUILT = -X '$(VERSION_PATH).BuildTime=$$(date "+%Y-%m-%d %H:%M:%S")'
STATUS = -X '$(VERSION_PATH).GitStatus=$$(git status)'
LDFLAG = "-s -w $(VERSION) $(GO_VERSION) $(GIT_COMMIT) $(BUILT) $(STATUS)"

PROJECT = $(shell pwd|grep -io "[^/]*$$")
# 二进制文件生成目录
BIN_DIR = ../../bin
$(warning   $(LDFLAG))
BUILD = go build -ldflags $(LDFLAG) -o $(BIN_DIR)/$(PROJECT) .

.PHONY: build
build: ## 编译为当前系统的二进制文件
	@$(BUILD)

.PHONY: run
run: ## 直接运行不编译
	@go run .

.PHONY: install
install: ## 安装二进制文件到系统path
	@chmod +x $(BIN_DIR)/$(PROJECT) && mv $(BIN_DIR)/$(PROJECT) /usr/local/bin

.PHONY: clean
clean: ## 清理
	@rm -rf bin	test	

.PHONY: docker
ifeq ($(REGISTRY),)
docker: ## 编译为docker镜像
	@docker build -t $(PROJECT):latest -t $(PROJECT):$(VER) -f Dockerfile .
else
docker: ## 编译为docker镜像 
	@docker build -t $(REGISTRY)/$(PROJECT):latest -t $(REGISTRY)/$(PROJECT):$(VER) -f Dockerfile .
endif

.PHONY: build_in_docker
build_in_docker: ## Dockerfile中执行编译
	@CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAG) .

.PHONY: windows
windows: ## 交叉编译为windows的二进制文件
	@GOOS=windows $(BUILD)

.PHONY: darwin
darwin: ## 交叉编译为苹果osx的二进制文件
	@GOOS=darwin $(BUILD)

.PHONY: arm
arm: ## 交叉编译为arm的linux环境（树莓派等环境）二进制文件
	@GOARCH=arm GOARM=7 GOOS=linux $(BUILD)

.PHONY: help
help: ## 显示make的目标
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)
