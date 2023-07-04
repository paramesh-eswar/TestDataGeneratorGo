package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func init() {
	fmt.Println("=======================")
	fmt.Println("Test Data Generator")
	fmt.Println("=======================")
}

var topWindow fyne.Window

func main() {
	a := app.New()
	tdgIcon, _ := fyne.LoadResourceFromPath("./theme/icons/tdg.png")
	// a.SetIcon(theme.TdgLogo())
	a.SetIcon(tdgIcon)
	// a.Settings().SetTheme(theme.MyTdgTheme{})
	logLifecycle(a)
	w := a.NewWindow("Test Data Generator")
	w.Resize(fyne.NewSize(600, 400))
	w.SetMaster()
	topWindow = w

	// title := canvas.NewText("Test Data Generator", color.Black)
	title := widget.NewLabel("Test Data Generator")
	title.TextStyle.Bold = true
	titleContainer := container.NewCenter(title)

	var fileContents string
	var fileName string
	var fileChooser *dialog.FileDialog

	numOfRowsLbl := widget.NewRichText(&widget.TextSegment{Text: "Number of rows", Style: widget.RichTextStyle{Alignment: fyne.TextAlignTrailing, ColorName: widget.RichTextStyleCodeInline.ColorName}})
	numOfRowsEntry := widget.NewEntry()
	numOfRowsEntry.SetPlaceHolder("Enter number of rows")

	metadataFileLbl := widget.NewRichText(&widget.TextSegment{Text: "Metadata file", Style: widget.RichTextStyle{Alignment: fyne.TextAlignTrailing, ColorName: widget.RichTextStyleCodeInline.ColorName}})
	metadataFileEntry := widget.NewEntry()
	metadataFileEntry.Disable()
	metadataFileEntry.SetPlaceHolder("Metadata file path")
	fileUploadBtn := widget.NewButton("Choose file", func() {
		fileChooser = dialog.NewFileOpen(func(r fyne.URIReadCloser, _ error) {
			data, _ := ioutil.ReadAll(r)
			result := fyne.NewStaticResource("name", data)
			fileContents = string(result.StaticContent)
			fileName = r.URI().Path()
			// log.Print(fileName)
			// log.Print(fileContents)
			// fw := fyne.CurrentApp().NewWindow(string(result.StaticName))
			// fw.SetContent(container.NewScroll(resultArea))
			// fw.Show()
		}, w)
		fileChooser.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileChooser.Show()
		log.Print(fileName)
		log.Print(fileContents)
	})
	// fileUploadBtn.Resize(fyne.Size{Width: 10, Height: 30})

	// title.Move(fyne.NewPos(20, 20))
	// content := container.NewBorder(titleContainer, nil, nil, nil, nil)
	content := container.New(layout.NewVBoxLayout(), titleContainer, container.New(layout.NewAdaptiveGridLayout(3), numOfRowsLbl, numOfRowsEntry, layout.NewSpacer(), metadataFileLbl, metadataFileEntry, fileUploadBtn))
	w.SetContent(content)
	w.SetFixedSize(true)
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}
