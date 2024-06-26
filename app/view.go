package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type View struct {
	store Store
}

func NewView(store Store) View {
	return View{store: store}
}

func (v View) Build() error {
	app := tview.NewApplication()

	allTodos, err := v.store.GetAll()
	if err != nil {
		return err
	}

	taskTable := tview.NewTable().SetBorders(true)
	taskTable.
		SetBorder(true).
		SetTitle("All tasks")

	priority := "High"

	textArea := tview.NewTextArea().
		SetPlaceholder("Describe your task...").
		SetPlaceholderStyle(tcell.Style{}.Dim(true))

	highPButton := tview.NewButton("High").
		SetSelectedFunc(func() {
			priority = "High"
		})

	mediumPButton := tview.NewButton("Medium").
		SetSelectedFunc(func() {
			priority = "Medium"
		})

	lowPButton := tview.NewButton("Low").SetSelectedFunc(func() {
		priority = "Low"
	})

	createButton := tview.NewButton("Add").
		SetStyle(tcell.Style{}.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite).Bold(true)).
		SetBackgroundColorActivated(tcell.ColorDarkGreen).
		SetLabelColorActivated(tcell.ColorWhite).
		SetSelectedFunc(func() {
			if strings.TrimSpace(textArea.GetText()) == "" {
				return
			}

			todo, dbErr := v.store.Add(textArea.GetText(), priority)
			allTodos = append([]Todo{todo}, allTodos...)

			renderTableRows(taskTable, allTodos, "All")

			textArea.SetText("", true)

			fmt.Printf("error creating todo %v", dbErr)
		})

	createTaskFlex := tview.NewFlex().
		AddItem(textArea, 0, 15, true).
		AddItem(tview.NewFlex().
			AddItem(highPButton, 0, 1, false).
			AddItem(mediumPButton, 0, 1, false).
			AddItem(lowPButton, 0, 1, false).
			SetDirection(tview.FlexColumn),
			0, 1, false).
		AddItem(createButton, 0, 1, false).
		SetDirection(tview.FlexRow)
	createTaskFlex.
		SetBorder(true).
		SetTitle("New")

	priorityList := tview.NewList().ShowSecondaryText(false).
		AddItem("All", "", 0, func() {
			renderTableRows(taskTable, allTodos, "All")
		}).
		AddItem("High", "", 0, func() {
			renderTableRows(taskTable, allTodos, "High")
		}).
		AddItem("Medium", "", 0, func() {
			renderTableRows(taskTable, allTodos, "Medium")
		}).
		AddItem("Low", "", 0, func() {
			renderTableRows(taskTable, allTodos, "Low")
		}).
		SetCurrentItem(0).
		SetHighlightFullLine(true).
		SetMainTextStyle(tcell.Style{}.Bold(true))
	priorityList.
		SetBorder(true).
		SetTitle("Filter")

	leftFlex := tview.NewFlex().
		AddItem(createTaskFlex, 0, 3, true).
		AddItem(priorityList, 0, 1, false).
		SetDirection(tview.FlexRow)

	taskTable.SetSelectedFunc(func(row, column int) {
		currentStatus := taskTable.GetCell(row, 2).Text
		var newStatus Status
		var color tcell.Color

		if currentStatus == string(StatusPending) {
			newStatus = StatusDone
			color = tcell.ColorGreen
		} else {
			newStatus = StatusPending
			color = tcell.ColorRed
		}
		task := allTodos[row-1]
		(&task).Status = newStatus
		allTodos[row-1] = task

		_ = v.store.Update(task.ID, newStatus)

		taskTable.GetCell(row, 2).SetText(string(newStatus)).SetTextColor(color)
	}).
		SetSelectable(true, false).
		Select(1, 0)

	renderTableRows(taskTable, allTodos, "All")

	usageDirection := tview.NewTextView()
	mouse := "Use your mouse to switch between different sections of the app"
	navigate := fmt.Sprintf("%s%s To navigate up/down the task list", string(tcell.RuneUArrow), string(tcell.RuneDArrow))
	toggle := fmt.Sprintf("Press Enter to toggle task between done/pending after selecting a row")
	quit := fmt.Sprintf("Press %s+%s to exit the app", "ctrl", "c")

	sb := strings.Builder{}
	sb.WriteString(mouse)
	sb.WriteString("\n")
	sb.WriteString(navigate)
	sb.WriteString("\n")
	sb.WriteString(toggle)
	sb.WriteString("\n")
	sb.WriteString(quit)

	usageDirection.SetText(sb.String())
	usageDirection.
		SetBorder(true).
		SetTitle("Usage details")

	rightFlex := tview.NewFlex().
		AddItem(taskTable, 0, 8, false).
		AddItem(usageDirection, 0, 2, false).
		SetDirection(tview.FlexRow)

	mainOuterFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, true).
		AddItem(rightFlex, 0, 4, false).
		SetDirection(tview.FlexColumn)

	if err := app.SetRoot(mainOuterFlex, true).SetFocus(mainOuterFlex).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

func renderTableRows(taskTable *tview.Table, todos []Todo, priority string) {
	taskTable.Clear()

	taskTable.SetCell(0, 0, &tview.TableCell{
		Text:          "Priority",
		Align:         tview.AlignCenter,
		Expansion:     1,
		Color:         tcell.ColorWhite,
		NotSelectable: true,
	})
	taskTable.SetCell(0, 1, &tview.TableCell{
		Text:          "Task",
		Align:         tview.AlignCenter,
		Expansion:     4,
		Color:         tcell.ColorWhite,
		NotSelectable: true,
	})
	taskTable.SetCell(0, 2, &tview.TableCell{
		Text:          "Status",
		Align:         tview.AlignCenter,
		Expansion:     1,
		Color:         tcell.ColorWhite,
		NotSelectable: true,
	})

	rowIndex := 1

	for _, todo := range todos {
		if priority != "All" && todo.Priority != priority {
			continue
		}

		color := tcell.ColorRed
		if todo.Status == StatusDone {
			color = tcell.ColorGreen
		}

		taskTable.SetCell(rowIndex, 0, &tview.TableCell{
			Text:      todo.Priority,
			Align:     tview.AlignCenter,
			Expansion: 1,
		})
		taskTable.SetCell(rowIndex, 1, &tview.TableCell{
			Text:      todo.Task,
			Align:     tview.AlignCenter,
			Expansion: 1,
		})
		taskTable.SetCell(rowIndex, 2, &tview.TableCell{
			Text:      string(todo.Status),
			Align:     tview.AlignCenter,
			Expansion: 1,
			Color:     color,
		})
		rowIndex++
	}
}
