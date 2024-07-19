package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	keyTimerLength = "focus.default"
)

func main() {
	a := app.NewWithID("xyz.andy.fomato")
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	w := a.NewWindow("Fomato Timer")

	focusTime := a.Preferences().IntWithFallback(keyTimerLength, 30*60)
	timer := widget.NewRichText()
	updateTime(timer, focusTime)

	less := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		if focusTime <= 5*60 { // min bound
			return
		}

		focusTime -= 60 * 5
		updateTime(timer, focusTime)
		a.Preferences().SetInt(keyTimerLength, focusTime)
	})
	more := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		focusTime += 60 * 5
		updateTime(timer, focusTime)
		a.Preferences().SetInt(keyTimerLength, focusTime)
	})
	timeRow := container.NewHBox(container.NewCenter(less),
		padTime(timer),
		container.NewCenter(more))

	focus := widget.NewButton("Focus", func() {
		startTimer(focusTime, "Focus", w.Canvas())
	})
	focus.Importance = widget.HighImportance
	slack := widget.NewButton("Break", func() {
		startTimer(5*60, "Break", w.Canvas())
	})
	content := container.NewCenter(container.NewVBox(timeRow,
		container.NewGridWithColumns(2, slack, focus)))
	w.SetContent(container.NewPadded(container.NewPadded(content)))
	w.ShowAndRun()
}

func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60

	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func padTime(t *widget.RichText) fyne.CanvasObject {
	pad := theme.Padding()

	return container.New(layout.NewCustomPaddedLayout(-3.5*pad, -2.5*pad, pad, pad), t)
}

func startTimer(remain int, name string, c fyne.Canvas) {
	ticker := widget.NewRichText()
	updateTime(ticker, remain)

	stop := widget.NewButton("Stop", nil)
	overlay := container.NewPadded(container.NewVBox(
		padTime(ticker),
		stop))

	p := widget.NewModalPopUp(overlay, c)
	stop.OnTapped = func() {
		remain = -1 // don't notify
		p.Hide()
	}
	go func() {
		for remain > 0 {
			updateTime(ticker, remain)

			remain--
			time.Sleep(time.Second)
		}

		if remain == 0 {
			fyne.CurrentApp().SendNotification(fyne.NewNotification(name+" done",
				"Your "+strings.ToLower(name)+" timer finished"))
		}
		p.Hide()
	}()
	p.Show()
}

func updateTime(out *widget.RichText, time int) {
	out.ParseMarkdown("# " + formatTimer(time))
	themeTimer(out, time)
}
