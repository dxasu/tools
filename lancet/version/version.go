package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

/*
VERSION    = $(shell git describe --tags --always)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
git_tag    = $(shell git describe --tags --abbrev=0)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)
git_branch = $(shell git show-ref | grep $(shell git show HEAD | sed -n 1p | cut -d " " -f 2) | sed 's|.* /(.*)|1|' | grep -v HEAD | sort | uniq | head -n 1)
*/
var (
	Version   string
	BuildTime string
	GitCommit string
	GitStatus string
)

func init() {
	args := os.Args
	if len(args) == 2 && args[1] == "--version" {
		fmt.Println("\033[;32mBuilder:\033[0m")
		fotmat := "%12s %-10s\n"
		fmt.Printf(fotmat, "Version:", Version)
		fmt.Printf(fotmat, "GitCommit:", GitCommit)
		fmt.Printf(fotmat, "BuildTime:", BuildTime)
		fmt.Printf(fotmat, "GoVersion:", runtime.Version())
		fmt.Printf(fotmat, "OS/Arch:", runtime.GOOS+"/"+runtime.GOARCH)

		if GitStatus != "" && strings.Contains(GitStatus, ".go") {
			fmt.Println(
				"\033[;31mWARNING: Build on unclean git status:\033[0m")
			fmt.Println(GitStatus)
			fmt.Println()
		}
		os.Exit(0)
	}
}
