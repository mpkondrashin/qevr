/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_folder.go

Pick destination folder
*/
package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageSource struct {
	BasePage
	folderEntry *widget.Entry
}

var _ Page = &PageSource{}

func (p *PageSource) Index() PageIndex {
	return pgSource
}

func (p *PageSource) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgFilter
}

func (p *PageSource) Content() fyne.CanvasObject {
	labelFolder := widget.NewLabel("Source CSV file:")
	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(p.wiz.config.Source)
	folderButton := widget.NewButton("Change...", func() {
		folderDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.URI().Path())
		}, p.wiz.win)
		folderDialog.Show()
	})
	return container.NewVBox(labelFolder,
		container.NewBorder(nil, nil, nil, folderButton, p.folderEntry)) // p.folderEntry, folderButton)
}

func (p *PageSource) AquireData(config *Config) error {
	p.wiz.csv = NewCSV(strings.TrimSpace(p.folderEntry.Text))
	err := p.wiz.csv.DetectIPCVE()
	if err != nil {
		return err
	}
	config.Source = p.folderEntry.Text
	log.Printf("%s: ip: %d, cve: %d", config.Source, p.wiz.csv.IPIndex, p.wiz.csv.CVEIndex)
	return nil
}
