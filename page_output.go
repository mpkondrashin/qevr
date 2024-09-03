/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_output.go

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
	fileNameLabel *widget.Label
	suffixEntry   *widget.Entry
	folderEntry   *widget.Entry
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
	p.suffixEntry.SetText(p.wiz.config.Output.File.Suffix)
	p.suffixEntry.OnChanged = func(_ string) {
		if p.fileNameLabel == nil {
			return
		}
		p.fileNameLabel.SetText(p.FileName())
		p.fileNameLabel.Refresh()
	}
	suffixItem := widget.NewFormItem("Suffix for filename:", p.suffixEntry)

	p.fileNameLabel = widget.NewLabel(p.FileName())
	fileNameItem := widget.NewFormItem("Output filename:", p.fileNameLabel)

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
		fileNameItem,
		suffixItem,
		folderItem,
	)
}

func (p *PageOutput) AquireData(config *Config) error {
	p.wiz.config.Output.File.Folder = strings.TrimSpace(p.folderEntry.Text)
	p.wiz.config.Output.File.Suffix = strings.TrimSpace(p.suffixEntry.Text)
	return p.wiz.model.SaveToFile(filepath.Join(p.folderEntry.Text, p.FileName()))
}

func (p *PageOutput) Suffix() string {
	if p.suffixEntry == nil {
		return p.wiz.config.Output.File.Suffix

	}
	return p.suffixEntry.Text
}

func (p *PageOutput) FileName() string {
	currentDate := time.Now().Format("20060102")
	suffix := p.Suffix()
	if suffix != "" {
		suffix = "_" + suffix
	}
	return fmt.Sprintf("qevr_%s%s.csv", currentDate, suffix)
}
