/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_finish.go

Final installer page
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PageFinish struct {
	BasePage
	//runCheck *widget.Check
}

var _ Page = &PageFinish{}

func (p *PageFinish) Index() PageIndex {
	return pgFinish
}

func (p *PageFinish) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgExit
}

func (p *PageFinish) Content() fyne.CanvasObject {
	text := "Thank you for using QeVR.\nYou can now close this window."
	//The configuration has been saved successfully.

	//var labels []fyne.CanvasObject //  *widget.Label
	//for _, l := range strings.Split(text, "\n") {
	//labels = append(labels, widget.NewLabel(l))
	//}
	//return container.NewVBox(labels...)
	return widget.NewLabel(text)
}

func (p *PageFinish) AquireData(config *Config) error {
	folder, err := ExecutableFolder()
	if err != nil {
		return err
	}
	return p.wiz.config.SaveConfig(folder)
}
