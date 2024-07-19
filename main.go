package main

import (
	"fmt"
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
	pad := theme.Padding()
	timeRow := container.NewHBox(container.NewCenter(less),
		container.New(layout.NewCustomPaddedLayout(-3.5*pad, -2.5*pad, pad, pad), timer),
		container.NewCenter(more))

	start := widget.NewButton("Start", func() {
		ticker := widget.NewRichText()
		updateTime(ticker, focusTime)
		remain := focusTime
		stop := widget.NewButton("Stop", nil)
		overlay := container.NewPadded(container.NewVBox(
			container.New(layout.NewCustomPaddedLayout(-3.5*pad, -2.5*pad, pad, pad), ticker),
			stop))

		p := widget.NewModalPopUp(overlay, w.Canvas())
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
				a.SendNotification(fyne.NewNotification("Focus done", "Your focus timer finished"))
			}
			p.Hide()
		}()
		p.Show()
	})
	content := container.NewCenter(container.NewVBox(timeRow, start))
	w.SetContent(container.NewPadded(container.NewPadded(content)))
	w.ShowAndRun()
}

func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60

	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func updateTime(out *widget.RichText, time int) {
	out.ParseMarkdown("# " + formatTimer(time))
	themeTimer(out, time)
}
