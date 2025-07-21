package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/dxasu/pure/rain"
	"github.com/dxasu/pure/storage"
)

// https://github.com/search?q=go-btfs&type=repositories&l=Go&s=stars&o=desc

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "-o" {
		ctx := context.Background()
		obj := ""
		if len(os.Args) >= 3 {
			obj = os.Args[2]
		}
		out, err := rain.Shell(ctx, `fzf --walker dir,follow,hidden --walker-root `+os.Getenv("GOPATH")+"/pkg/mod/github.com/"+obj)
		rain.ExitIf(err)
		_, err = rain.Shell(ctx, "code "+out)
		rain.ExitIf(err)
		return
	}
	d := storage.NewDataPath(storage.DirTemp, "web_search").GetFilePath(os.Args[0] + ".url")
	if len(os.Args) >= 2 && os.Args[1] == "-u" {
		if len(os.Args) < 3 {
			fmt.Println("format file: " + d)
			return
		}
		r, err := url.Parse(os.Args[2])
		rain.ExitIf(err)
		if !r.IsAbs() {
			fmt.Println("Please check url format.")
			return
		}
		if strings.Count(r.String(), "%s") == 0 {
			rain.ExitIf("invalid url format, miss %s")
		}
		err = os.WriteFile(d, []byte(r.String()), 0755)
		rain.ExitIf(err)
		fmt.Println("Update url format successfully. You can use it later.")
		return
	}
	c, _ := os.ReadFile(d)
	if len(c) == 0 || len(os.Args) == 1 {
		fmt.Println("Usage: " + os.Args[0] + " -u <url format>")
		fmt.Println(os.Args[0] + " -u for look up the url format file path.")
		fmt.Println(os.Args[0] + " [https://]github.com for opening url")
		fmt.Println("Example: " + os.Args[0] + " -u https://github.com/search?q=%s&type=repositories&l=Go&s=stars&o=desc")
		fmt.Println("This will save to " + d)
		return
	}

	if strings.HasPrefix(os.Args[1], "github.com") || strings.HasPrefix(os.Args[1], "https://github.com") {
		u, err := url.Parse(os.Args[1])
		if err != nil {
			panic(err)
		}
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		err = rain.OpenBrower(u.String())
		rain.ExitIf(err)
		return
	}

	urlFormat := strings.TrimSpace(string(c))
	placeholderNum := strings.Count(urlFormat, "%s")
	var str rain.Clog = 0xff0000
	if len(os.Args) != 1+placeholderNum {
		fmt.Println("Usage: "+os.Args[0]+" xxx", placeholderNum, "placeholders required")
		fmt.Println("Current url need", str.Str(placeholderNum), "placeholder, format: "+urlFormat)
		return
	}

	anyArgs := make([]any, 0, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		anyArgs = append(anyArgs, os.Args[i])
	}
	u, err := url.Parse(fmt.Sprintf(urlFormat, anyArgs...))
	if err != nil {
		panic(err)
	}
	fmt.Println(u.String())
	err = rain.OpenBrower(u.String())
	rain.ExitIf(err)
}
