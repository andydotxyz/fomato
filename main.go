package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("xyz.andy.fomato")
	w := a.NewWindow("Fomato Timer")

	focusTime := 30 * 60
	timer := widget.NewRichTextFromMarkdown(formatTimerMarkdown(focusTime))

	start := widget.NewButton("Start", func() {
		timer := widget.NewRichTextFromMarkdown(formatTimerMarkdown(focusTime))
		remain := 30 * 60
		stop := widget.NewButton("Stop", nil)
		overlay := container.NewPadded(container.NewVBox(timer, stop))

		p := widget.NewModalPopUp(overlay, w.Canvas())
		stop.OnTapped = func() {
			remain = 0
			p.Hide()
		}
		go func() {
			for remain > 0 {
				timer.ParseMarkdown(formatTimerMarkdown(remain))

				remain--
				time.Sleep(time.Second)
			}

			p.Hide()
		}()
		p.Show()
	})
	content := container.NewCenter(container.NewVBox(timer, start))
	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()
}

func formatTimer(time int) string {
	secs := time % 60
	mins := (time - secs) / 60

	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func formatTimerMarkdown(sec int) string {
	return "# " + formatTimer(sec)
}
