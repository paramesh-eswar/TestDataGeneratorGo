package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tdgTheme "github.com/paramesh-eswar/TestDataGeneratorGo/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
	fmt.Println("=======================")
	fmt.Println("Test Data Generator")
	fmt.Println("=======================")
}

var topWindow fyne.Window
var numOfRows int
var metadataFileName string
var progressBar *widget.ProgressBar

type numEntry struct {
	widget.Entry
}

var lightThemeSelected bool = false
var PercentageCompleted = 0.0

// var externalFloat binding.ExternalFloat

func main() {
	a := app.New()
	tdgIcon := tdgTheme.TdgLogo()
	a.SetIcon(tdgIcon)
	a.Settings().SetTheme(tdgTheme.MyTdgDarkTheme{})
	// logLifecycle(a)
	w := a.NewWindow("Test Data Generator")
	w.Resize(fyne.NewSize(600, 400))
	w.SetMaster()
	topWindow = w

	// widget button to change the theme
	tdgThemeBtn := widget.NewButtonWithIcon("", tdgTheme.LightThemeIcon(), nil)
	tdgThemeBtn.OnTapped = func() {
		if lightThemeSelected {
			a.Settings().SetTheme(tdgTheme.MyTdgDarkTheme{})
			lightThemeSelected = false
			tdgThemeBtn.SetIcon(tdgTheme.LightThemeIcon())
		} else {
			a.Settings().SetTheme(tdgTheme.MyTdgLightTheme{})
			lightThemeSelected = true
			tdgThemeBtn.SetIcon(tdgTheme.DarkThemeIcon())
		}
	}

	title := widget.NewRichText(&widget.TextSegment{
		Text: "Test Data Generator",
		Style: widget.RichTextStyle{
			Alignment: fyne.TextAlignCenter,
			ColorName: widget.RichTextStyleCodeInline.ColorName,
			TextStyle: fyne.TextStyle{Bold: true},
			SizeName:  fyne.ThemeSizeName(theme.SizeNameHeadingText),
		},
	})
	titleContainer := container.NewCenter(title)

	var fileChooser *dialog.FileDialog

	numOfRowsLbl := widget.NewRichText(&widget.TextSegment{
		Text: "Number of rows",
		Style: widget.RichTextStyle{
			Alignment: fyne.TextAlignTrailing,
			ColorName: widget.RichTextStyleCodeInline.ColorName,
		},
	})
	numOfRowsEntry := &numEntry{}
	numOfRowsEntry.ExtendBaseWidget(numOfRowsEntry)
	numOfRowsEntry.SetPlaceHolder("Enter number of rows")
	numOfRowsEntry.Validator = validation.NewRegexp(`^[1-9][0-9]*$`, "Number of rows must be greater than zero")

	metadataFileLbl := widget.NewRichText(&widget.TextSegment{
		Text: "Metadata file",
		Style: widget.RichTextStyle{
			Alignment: fyne.TextAlignTrailing,
			ColorName: widget.RichTextStyleCodeInline.ColorName,
		},
	})
	metadataFileEntry := widget.NewEntry()
	metadataFileEntry.Disable()
	metadataFileEntry.SetPlaceHolder("Metadata file path")
	fileChooser = dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		// data, _ := ioutil.ReadAll(r)
		// result := fyne.NewStaticResource("name", data)
		// fileContents := string(result.StaticContent)
		if err != nil {
			log.Print("Error occured while choosing the metadata file\n")
			log.Print(err)
			return
		}

		if r == nil {
			log.Print("cancelled")
			return
		}

		fileName := r.URI().Path()
		// log.Print("inside:" + fileName)
		metadataFileEntry.SetText(fileName)
		persistFileName(fileName)
	}, w)
	fileChooser.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fileChooser.Resize(fyne.NewSize(500, 400))

	fileUploadBtn := widget.NewButton("Choose file", func() {
		fileChooser.Show()
	})

	resultPane := widget.NewMultiLineEntry()
	resultPane.SetMinRowsVisible(4)
	resultPane.Wrapping = fyne.TextWrapBreak
	resultPane.Disable()

	// progress bar
	// externalFloat = binding.BindFloat(&PercentageCompleted)
	// progressBar = widget.NewProgressBarWithData(externalFloat)
	progressBar = widget.NewProgressBar()

	generateDataBtn := widget.NewButton("Generate Data", func() {
		// PercentageCompleted = 0.0
		// externalFloat.Reload()
		resultPane.SetText("")

		if numOfRowsEntry.Validate() != nil {
			dialog.ShowError(errors.New(numOfRowsEntry.Validate().Error()), topWindow)
			topWindow.Canvas().Focus(numOfRowsEntry)
			return
		}
		if len(metadataFileName) == 0 || !strings.HasSuffix(metadataFileName, ".json") {
			dialog.ShowError(errors.New("Invalid metadata file path found"), topWindow)
			return
		}

		numOfRowsInt, err := strconv.Atoi(numOfRowsEntry.Text)
		if err != nil {
			resultPane.SetText(err.Error())
		}

		numOfRows = numOfRowsInt

		output := validateMetadata()
		if strings.EqualFold(output, "success") {
			output = testDataGenerator()
		}
		resultPane.SetText("Output from test data generator:\n" + output)
	})

	content := container.NewVBox(
		container.NewGridWithColumns(3,
			layout.NewSpacer(), titleContainer, container.NewHBox(
				layout.NewSpacer(),
				container.NewGridWrap(fyne.NewSize(50, 40), tdgThemeBtn),
			),
		),
		container.NewGridWithColumns(3,
			numOfRowsLbl, numOfRowsEntry, layout.NewSpacer(),
			metadataFileLbl, metadataFileEntry, container.NewGridWrap(fyne.NewSize(150, 40), fileUploadBtn),
			layout.NewSpacer(), generateDataBtn, layout.NewSpacer(),
		),
		container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 30)), progressBar)),
		container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 200)), resultPane)),
	)
	w.SetContent(content)
	w.SetFixedSize(true)
	w.ShowAndRun()
}

func persistFileName(fileName string) {
	metadataFileName = fileName
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
