/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

Intro page
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
	IntoText = "QeVR helps to upload vulnerability scans to Tipping Point SMS server for profile tuning."
)

func (p *PageIntro) Content() fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle("QeVR",
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	version := "Version 1.0" //fmt.Sprintf("Version %s build %s", globals.Version, globals.Build)
	versionLabel := widget.NewLabelWithStyle(version,
		fyne.TextAlignCenter, fyne.TextStyle{})

	report := widget.NewRichTextFromMarkdown(IntoText)
	report.Wrapping = fyne.TextWrapWord

	repoURL, _ := url.Parse("https://github.com/mpkondrashin/qevr")
	repoLink := widget.NewHyperlink("QeVR repository on GitHub", repoURL)

	licensePopUp := func() {
		licenseLabel := widget.NewLabel(LicenseText())
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
		container.NewHBox(repoLink, licenseButton),
	)
}

func (p *PageIntro) AquireData(config *Config) error {
	// check acept license
	return nil
}

func LicenseText() string {
	return `MIT License

Copyright (c) 2024 Michael Kondrashin (mkondrashin@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
}
