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
	rowChar     = "" // -
	columnChar  = "" // '|'
	minify      = true
	header      = true
	braceFormat = ""
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

		if braceFormat != "" {
			r := strings.NewReplacer("\\n", "\n", "\\t", "\t")
			// "f,f"f(f,f)
			braceFormat = r.Replace(braceFormat)
			holds := strings.Split(braceFormat, "f")
			if len(holds) < 6 {
				panic(`braceFormat must like: "f,f"f(f,f)`)
			}
			arr := ToArray(data)
			output := strings.Builder{}
			for i, row := range arr {
				if i != 0 {
					output.WriteString(holds[4])
				}
				output.WriteString(holds[3])
				for j, v := range row {
					if j != 0 {
						output.WriteString(holds[1])
					}
					output.WriteString(holds[0])
					output.WriteString(v)
					output.WriteString(holds[2])
				}
				output.WriteString(holds[5])
			}
			clipboard.WriteAll(output.String())
			return
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
	rootCmd.Flags().StringVarP(&braceFormat, "brace format", "f", "", "brace")
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if helpFlag, _ := cmd.Flags().GetBool("help"); helpFlag {
			cmd.Usage()
		}
	})
}

func ToArray(data any) [][]string {
	var rows [][]string
	switch data := data.(type) {
	case string:
		d := strings.Split(data, "\n")
		for _, line := range d {
			rows = append(rows, strings.Split(line, "\t"))
		}
	case []string:
		for _, line := range data {
			rows = append(rows, strings.Split(line, "\t"))
		}
	case [][]string:
		rows = data
	default:
	}
	return rows
}
