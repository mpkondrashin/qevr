/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageIntro struct {
	BasePage
}

var _ Page = &PageIntro{}

func (p *PageIntro) Index() PageIndex {
	return pgIntro
}

func (p *PageIntro) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgSource

}

const (
	IntoText = "QeVR helps to upload Qualys security scans to Tipping Point SMS server for profile tuning."
)

func (p *PageIntro) Content() fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle("QeVR",
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	version := "123" //fmt.Sprintf("Version %s build %s", globals.Version, globals.Build)
	versionLabel := widget.NewLabelWithStyle(version,
		fyne.TextAlignCenter, fyne.TextStyle{})

	report := widget.NewRichTextFromMarkdown(IntoText)
	report.Wrapping = fyne.TextWrapWord

	//noteMarkdown := widget.NewRichTextFromMarkdown("NoteText")
	//noteMarkdown.Wrapping = fyne.TextWrapWord

	repoURL, _ := url.Parse("https://github.com/mpkondrashin/qevr")
	repoLink := widget.NewHyperlink("QeVR repository on GitHub", repoURL)

	licensePopUp := func() {
		licenseLabel := widget.NewLabel("license") //LicenseText())
		sc := container.NewScroll(licenseLabel)
		popup := dialog.NewCustom("Show License Information", "Close", sc, p.wiz.win)
		popup.Resize(fyne.NewSize(800, 600))
		popup.Show()
	}
	licenseButton := widget.NewButton("License Information...", licensePopUp)
	return container.NewVBox(
		titleLabel,
		versionLabel,
		report,
		//noteMarkdown,
		container.NewHBox(repoLink, licenseButton),
	)
}

func (p *PageIntro) Run() {
	//	fmt.Println("Run" + p.Name())
	//	fmt.Println("Type ", p.wiz.installer.config.GetString(config.Engine.String()))
	//
	// p.sandboxRadio.SetSelected(p.wiz.installer.config.SandboxType.String())
}

func (p *PageIntro) AquireData(config *Config) error {
	// check acept license
	return nil
}

/*
func LicenseText() string {
	filePath := "embed/LICENSE"
	licFile, err := embedFS.Open(filePath)
	if err != nil {
		return "reading error"
	}
	defer func() {
		licFile.Close()
	}()
	licBytes, err := io.ReadAll(licFile)
	if err != nil {
		return "reading error"
	}
	return string(licBytes)

}
*/
