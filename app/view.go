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
			_, dbErr := v.store.Add(textArea.GetText(), priority)
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

	taskList := tview.NewTable().SetBorders(true)
	taskList.
		SetBorder(true).
		SetTitle("All tasks")

	taskList.SetCell(0, 0, &tview.TableCell{
		Text:      "Priority",
		Align:     tview.AlignCenter,
		Expansion: 1,
		Color:     tcell.ColorBlue,
	})
	taskList.SetCell(0, 1, &tview.TableCell{
		Text:      "Task",
		Align:     tview.AlignCenter,
		Expansion: 4,
		Color:     tcell.ColorBlue,
	})
	taskList.SetCell(0, 2, &tview.TableCell{
		Text:      "Status",
		Align:     tview.AlignCenter,
		Expansion: 1,
		Color:     tcell.ColorBlue,
	})

	allTodos, err := v.store.GetAll()
	if err != nil {
		return err
	}

	for index, todo := range allTodos {
		taskList.SetCell(index+1, 0, &tview.TableCell{
			Text:      todo.Priority,
			Align:     tview.AlignCenter,
			Expansion: 1,
			Color:     tcell.ColorBlue,
		})
		taskList.SetCell(index+1, 1, &tview.TableCell{
			Text:      todo.Task,
			Align:     tview.AlignCenter,
			Expansion: 1,
			Color:     tcell.ColorBlue,
		})
		taskList.SetCell(index+1, 2, &tview.TableCell{
			Text:      string(todo.Status),
			Align:     tview.AlignCenter,
			Expansion: 1,
			Color:     tcell.ColorBlue,
		})
	}

	mainOuterFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, true).
		AddItem(taskList, 0, 4, false).
		SetDirection(tview.FlexColumn)

	if err := app.SetRoot(mainOuterFlex, true).SetFocus(mainOuterFlex).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}
