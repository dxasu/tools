package main

import (
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dxasu/pure/storage"
)

type menuItem struct {
	label   string
	command string
	output  string
}

func runCommand(item menuItem) func() {
	return func() {
		parts := strings.Fields(item.command)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			item.output = err.Error()
		}
	}
}

func main() {
	config, err := storage.New("config.json", true)
	if err != nil {
		panic(err)
	}
	const listKey = "buttons"
	if len(os.Args) > 3 && os.Args[1] == "set" {
		config.Set(listKey+"."+os.Args[2], strings.Join(os.Args[3:], " "))
		return
	}

	var menuItems []menuItem
	config.Range(listKey, func(k string, v any) {
		menuItems = append(menuItems, menuItem{k, v.(string), ""})
	})

	a := app.New()
	w := a.NewWindow("ctool")

	buttons := make([]fyne.CanvasObject, len(menuItems))
	for i, item := range menuItems {
		if item.label == "exit" {
			buttons[i] = widget.NewButton("exit", func() {
				w.Close()
			})
			continue
		}
		buttons[i] = widget.NewButton(item.label, runCommand(item))
	}

	content := container.NewVBox(buttons...)
	w.SetContent(content)
	minHeight := 1
	w.Resize(fyne.NewSize(150, float32(minHeight)))
	w.ShowAndRun()
}
