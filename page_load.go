/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_upload.go

Read CSV file
*/
package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageLoad struct {
	BasePage
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
}

var _ Page = &PageLoad{}

func (p *PageLoad) Index() PageIndex {
	return pgLoad
}

func (p *PageLoad) Name() string {
	return "Copy Files"
}

func (p *PageLoad) Next(previousPage PageIndex) PageIndex {
	p.previousPage = previousPage
	return pgTarget
}

func (p *PageLoad) Content() fyne.CanvasObject {
	colsLabel := widget.NewLabel(p.wiz.csv.String())
	p.progressBar = widget.NewProgressBar()
	p.statusLabel = widget.NewLabel("")
	return container.NewVBox(
		colsLabel,
		p.progressBar,
		p.statusLabel,
	)
}

func (p *PageLoad) Run() {
	p.statusLabel.SetText("Loading...")
	size, err := FileSize(p.wiz.config.Source)
	if err != nil {
		p.statusLabel.SetText("Failed")
		err = fmt.Errorf("%s: %w", "Failed:", err)
		if err != nil {
			log.Println(err)
		}
		dialog.ShowError(err, p.wiz.win)
		return
	}
	var current int64 = 0
	err = p.wiz.csv.Load(
		func(i int64) {
			current += i
			fraction := float64(current) / float64(size)
			p.progressBar.SetValue(fraction)
			p.statusLabel.SetText(p.wiz.model.Status())
		},
		func(ip, cve string) {
			p.wiz.model.AddIPAndCVE(ip, cve)
		},
	)
	p.statusLabel.SetText(p.wiz.model.FinalStatus())
	if err != nil {
		p.statusLabel.SetText("Failed")
		err = fmt.Errorf("%s: %w", "Failed:", err)
		if err != nil {
			log.Println(err)
		}
		dialog.ShowError(err, p.wiz.win)
		return
	}
}

func FileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
