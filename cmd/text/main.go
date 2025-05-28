package main

import (
	"bytes"
	"fmt"
	"os"
	"text/array"
	"text/csv"

	"strings"

	"github.com/atotto/clipboard"
	"github.com/dxasu/pure/rain"
	"github.com/dxasu/pure/stdin"
	"github.com/dxasu/pure/text"
	_ "github.com/dxasu/pure/version"
	"github.com/spf13/cobra"
)

var (
	rowChar    = "" // -
	columnChar = "" // '|'
	minify     = true
	header     = true
)

// 根命令
var rootCmd = &cobra.Command{
	Use:   "text",
	Short: "text formatting",
	Long:  `text is a command line tool for text formatting`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 主命令逻辑（若直接运行根命令时触发）
	},
	Run: func(cmd *cobra.Command, args []string) {
		var data string
		if len(args) > 0 {
			data = strings.Join(args, " ")
		} else {
			_data, err := stdin.GetStdin()
			if err == stdin.ErrNoStdin {
				err = nil
			}
			rain.ExitIf(err)
			data = string(_data)
			if len(data) == 0 {
				data, _ = clipboard.ReadAll()
			}
		}

		output := bytes.NewBuffer(nil)
		t := text.NewText(output, header, data)
		t.SetSymbols(&text.SymbolCustom{
			Row:    rowChar,
			Column: columnChar,
		})
		t.Flush()
		if !minify {
			fmt.Println(output.String())
			return
		}

		lines := bytes.Split(output.Bytes(), []byte("\n"))
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			fmt.Println(string(line))
		}
	},
}

func main() {
	csv.Init(rootCmd)
	array.Init(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rowChar, "row", "r", "", "row character")
	rootCmd.Flags().StringVarP(&columnChar, "column", "c", "", "column character")
	rootCmd.Flags().BoolVarP(&minify, "minify", "m", true, "minify output when char is empty")
	rootCmd.Flags().BoolVarP(&header, "header", "", true, "need header")
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if helpFlag, _ := cmd.Flags().GetBool("help"); helpFlag {
			cmd.Usage()
		}
	})
}
