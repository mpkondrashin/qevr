/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_target.go

Pick target: SMS CSV csv file
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageTarget struct {
	BasePage
	choice *widget.RadioGroup
}

var _ Page = &PageTarget{}

func (p *PageTarget) Index() PageIndex {
	return pgTarget
}

func (p *PageTarget) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if p.choice == nil {
		return pgSave
	}
	switch MapTargetLabelFromString[p.choice.Selected] {
	case TargetFile:
		return pgSave
	case TargetSMS:
		return pgSMS
	}
	return pgSave
}

func (p *PageTarget) Content() fyne.CanvasObject {
	label := widget.NewLabel(p.wiz.model.FinalStatus())
	p.choice = widget.NewRadioGroup(TargetLabels, p.Choice)
	p.choice.SetSelected(TargetLabels[p.wiz.config.Output.Target])
	return container.NewVBox(label, p.choice)
}

func (p *PageTarget) Choice(chosen string) {
	p.wiz.UpdatePagesList()
}

//ParseCIDR: takes a string representing an IP/mask and returns an IP and an IPNet
//IPNet.Contains: c

func (p *PageTarget) AquireData(config *Config) error {
	config.Output.Target = MapTargetLabelFromString[p.choice.Selected]
	return nil
}
