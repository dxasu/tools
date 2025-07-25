package csv

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/spf13/cobra"
)

var csvCmd = &cobra.Command{
	Use:   "table",
	Short: "table formatting",
	Long:  `table is a command line tool for table formatting`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 主命令逻辑（若直接运行根命令时触发）
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// 全局标志（对所有子命令生效）
func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(csvCmd)
}

func markdown(cmd *cobra.Command, args []string) {
	var buf bytes.Buffer
	table := tablewriter.NewTable(&buf,
		tablewriter.WithRenderer(renderer.NewMarkdown()),
	)
	table.Header([]string{"Name", "Age", "City"})
	table.Append([]string{"Alice", "25", "New York"})
	table.Append([]string{"Bob", "30", "Boston"})
	table.Render()
}

func html(cmd *cobra.Command, args []string) {
	var buf bytes.Buffer
	table := tablewriter.NewTable(&buf,
		tablewriter.WithRenderer(renderer.NewHTML()),
	)
	table.Header([]string{"Name", "Age", "City"})
	table.Append([]string{"Alice", "25", "New York"})
	table.Append([]string{"Bob", "30", "Boston"})
	table.Render()
}

func csv(cmd *cobra.Command, args []string) {
	var buf bytes.Buffer
	table := tablewriter.NewTable(&buf) // tablewriter.WithRenderer(renderer),

	table.Header([]string{"Name", "Age", "City"})
	table.Append([]string{"Alice", "25", "New York"})
	table.Append([]string{"Bob", "30", "Boston"})
	table.Render()
}
