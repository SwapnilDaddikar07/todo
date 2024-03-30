package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			todo, dbErr := v.store.Add(textArea.GetText(), priority)
			allTodos = append(allTodos, todo)

			taskTable.SetCell(len(allTodos), 0, &tview.TableCell{
				Text:      todo.Priority,
				Align:     tview.AlignCenter,
				Expansion: 1,
			})
			taskTable.SetCell(len(allTodos), 1, &tview.TableCell{
				Text:      todo.Task,
				Align:     tview.AlignCenter,
				Expansion: 1,
			})
			taskTable.SetCell(len(allTodos), 2, &tview.TableCell{
				Text:      string(todo.Status),
				Align:     tview.AlignCenter,
				Expansion: 1,
				Color:     tcell.ColorRed,
			})

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
		AddItem("All", "", 0, nil).
		AddItem("High", "", 0, nil).
		AddItem("Medium", "", 0, nil).
		AddItem("Low", "", 0, nil).
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
	}).SetSelectedFunc(func(row, column int) {
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
		_ = v.store.Update(task.ID, newStatus)

		taskTable.GetCell(row, 2).SetText(string(newStatus)).SetTextColor(color)
	}).
		SetSelectable(true, false).
		Select(1, 0)

	for index, todo := range allTodos {
		color := tcell.ColorRed
		if todo.Status == StatusDone {
			color = tcell.ColorGreen
		}

		taskTable.SetCell(index+1, 0, &tview.TableCell{
			Text:      todo.Priority,
			Align:     tview.AlignCenter,
			Expansion: 1,
		})
		taskTable.SetCell(index+1, 1, &tview.TableCell{
			Text:      todo.Task,
			Align:     tview.AlignCenter,
			Expansion: 1,
		})
		taskTable.SetCell(index+1, 2, &tview.TableCell{
			Text:      string(todo.Status),
			Align:     tview.AlignCenter,
			Expansion: 1,
			Color:     color,
		})
	}

	mainOuterFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, true).
		AddItem(taskTable, 0, 4, false).
		SetDirection(tview.FlexColumn)

	if err := app.SetRoot(mainOuterFlex, true).SetFocus(mainOuterFlex).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}
