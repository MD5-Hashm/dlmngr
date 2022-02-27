package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	dl "hashm.tech/dlmngr/downloader"
)

var (
	buttons   = []string{"ADD", "HTTP"}
	downloads []chan dl.Bw
	curdown   []string
)

func unblock(channel chan dl.Bw) {
	for {
		<-channel
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func f() {
	a := app.New()
	w := a.NewWindow("dlmngr")
	tt := widget.NewLabel("Please select an option from the sidebar")

	listview := widget.NewList(func() int {
		return len(buttons)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(buttons[id])
	})

	listview.OnSelected = func(id widget.ListItemID) {
		if id == 0 {
			input := widget.NewEntry()
			input.SetPlaceHolder("HTTP")
			addCont := container.NewHSplit(listview, container.NewVBox(input, widget.NewButton("GET", func() {

				if !stringInSlice(input.Text, curdown) {
					curdown = append(curdown, input.Text)
					out := make(chan dl.Bw)

					go dl.DownloadFile("", input.Text, out)
					go unblock(out)
					downloads = append(downloads, out)

					dllist := widget.NewList(func() int {
						return len(downloads)
					}, func() fyne.CanvasObject {
						return widget.NewLabel("")
					}, func(id widget.ListItemID, object fyne.CanvasObject) {
						out := downloads[id]
						name := (<-out).Name
						object.(*widget.Label).SetText(name)
					})
					cont := container.NewHSplit(listview, dllist)
					cont.Offset = 0.18

					dllist.OnSelected = func(id widget.ListItemID) {
						out := downloads[id]
						dlstr := strconv.FormatInt((<-out).Bwr, 10)
						cv := container.NewVBox(widget.NewLabel(dlstr))
						cont = container.NewHSplit(listview, cv)
					}

					w.SetContent(cont)
				} else {
					dialog.ShowConfirm("Already downloaded this file", "Already downloaded this file, Are you sure you want to redownload it?", func(tf bool) {
						if tf {
							curdown = append(curdown, input.Text)
							out := make(chan dl.Bw)
							go dl.DownloadFile("", input.Text, out)
							go unblock(out)
							downloads = append(downloads, out)

							dllist := widget.NewList(func() int {
								return len(downloads)
							}, func() fyne.CanvasObject {
								return widget.NewLabel("")
							}, func(id widget.ListItemID, object fyne.CanvasObject) {
								out := downloads[id]
								name := (<-out).Name
								object.(*widget.Label).SetText(name)
							})
							cont := container.NewHSplit(listview, dllist)
							cont.Offset = 0.18

							dllist.OnSelected = func(id widget.ListItemID) {
								out := downloads[id]
								dlstr := strconv.FormatInt((<-out).Bwr, 10)
								cv := container.NewVBox(widget.NewLabel(dlstr))
								cont = container.NewHSplit(listview, cv)
							}

							w.SetContent(cont)
						}
					}, w)
				}
			})))
			addCont.Offset = 0.18
			w.SetContent(addCont)
		} else if id == 1 {
			dllist := widget.NewList(func() int {
				return len(downloads)
			}, func() fyne.CanvasObject {
				return widget.NewLabel("")
			}, func(id widget.ListItemID, object fyne.CanvasObject) {
				out := downloads[id]
				name := (<-out).Name
				object.(*widget.Label).SetText(name)
			})
			cont := container.NewHSplit(listview, dllist)
			cont.Offset = 0.18
			w.SetContent(cont)
		}
	}
	split := container.NewHSplit(
		listview,
		tt)
	split.Offset = 0.18

	w.SetContent(split)
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}

func main() {
	f()
}
