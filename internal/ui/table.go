package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

func PaintTable(rows [][]string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "AWS CLOUDFORMATION DEBUGGER"
	p.Text = "PRESS q TO NEXT ERROR"
	p.SetRect(0, 0, 50, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	table := widgets.NewTable()
	table.Title = "Cloudformation Stack Failed"
	table.Rows = rows
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.BorderStyle = ui.NewStyle(ui.ColorRed)
	table.SetRect(0, 6, 200, 20)
	table.FillRow = true
	table.ColumnWidths = []int{20, 180}

	ui.Render(p, table)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}