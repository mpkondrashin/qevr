/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_folder.go

Pick destination folder
*/
package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageOutput struct {
	BasePage
	suffixEntry *widget.Entry
	folderEntry *widget.Entry
}

var _ Page = &PageOutput{}

func (p *PageOutput) Index() PageIndex {
	return pgSave
}

func (p *PageOutput) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgFinish
}

func (p *PageOutput) Content() fyne.CanvasObject {
	p.suffixEntry = widget.NewEntry()
	suffixItem := widget.NewFormItem("Suffix for filename:", p.suffixEntry)

	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(p.wiz.config.Output.File.Folder)
	folderButton := widget.NewButton("Change...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.Path())
		}, p.wiz.win)
		folderDialog.Show()
	})
	folderItem := widget.NewFormItem("Save to folder:",
		container.NewBorder(nil, nil, nil, folderButton, p.folderEntry),
	)
	return widget.NewForm(
		suffixItem,
		folderItem,
	)
	/*suffixLabel := widget.NewLabel("Add suffix to file name:")
	p.suffixEntry = widget.NewEntry()
	suffixLine := container.NewHBox(suffixLabel, p.suffixEntry)

	labelFolder := widget.NewLabel("Save to folder:")
	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(p.wiz.config.Output.File.Folder)
	folderButton := widget.NewButton("Change...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.Path())
		}, p.wiz.win)
		folderDialog.Show()
	})
	return container.NewVBox(
		suffixLine,
		labelFolder,
		container.NewBorder(nil, nil, nil, folderButton, p.folderEntry),
	)*/
}

func (p *PageOutput) AquireData(config *Config) error {
	p.wiz.config.Output.File.Folder = strings.TrimSpace(p.folderEntry.Text)
	p.wiz.config.Output.File.Prefix = strings.TrimSpace(p.suffixEntry.Text)
	currentDate := time.Now().Format("20060102")
	suffix := p.suffixEntry.Text
	if suffix != "" {
		suffix = "_" + suffix
	}
	fileName := filepath.Join(p.folderEntry.Text, fmt.Sprintf("qevr_%s%s.csv", currentDate, suffix))
	return p.wiz.model.SaveToFile(fileName)
}
