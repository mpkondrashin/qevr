/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_token.go

Provide Vision One token and domain
*/
package main

import (
	"fmt"
	"log"

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
	//fmt.Println("XXX ", p.choice.Selected, MapTargetFromString[p.choice.Selected], TargetFile, TargetSMS)
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
	fmt.Println("PageTarget>Content:", p.wiz.config.Output.Target, "label:", TargetLabels[p.wiz.config.Output.Target])
	p.choice.SetSelected(TargetLabels[p.wiz.config.Output.Target])
	return container.NewVBox(label, p.choice)
}

func (p *PageTarget) Choice(chosen string) {
	log.Println("PageFilter->Choice: ", chosen)
	p.wiz.UpdatePagesList()
}

func (p *PageTarget) Run() {
	// No need to load, config is loaded when application started
	//	err := installer.LoadConfig()
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, win)
	//	}
	//fmt.Println("Run" + p.Name())
}

//ParseCIDR: takes a string representing an IP/mask and returns an IP and an IPNet
//IPNet.Contains: c

func (p *PageTarget) AquireData(config *Config) error {
	config.Output.Target = MapTargetLabelFromString[p.choice.Selected]
	//fmt.Println("PageTarget) AquireData:", p.wiz.config.Output.Target, "label:", TargetLabels[p.wiz.config.Output.Target])
	return nil
}
