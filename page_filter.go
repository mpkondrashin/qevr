/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_filer.go

Pick filters
*/
package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageFilter struct {
	BasePage
	choice *widget.RadioGroup
	list   *widget.Entry
}

var _ Page = &PageFilter{}

func (p *PageFilter) Index() PageIndex {
	return pgFilter
}

func (p *PageFilter) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgLoad
}

func (p *PageFilter) Content() fyne.CanvasObject {
	p.choice = widget.NewRadioGroup(FilterTypesLabels, p.Choice)
	p.choice.SetSelected(FilterTypesLabels[p.wiz.config.Filter.Type])
	p.list = widget.NewMultiLineEntry()
	p.list.SetText(strings.Join(p.wiz.config.Filter.Networks, "\n"))
	hintLabel := widget.NewLabel(`Specify networks in form "2.1.0.0/24" (without quites).
Put one network per line.
"#" character can be used for comment.`)
	return container.NewVBox(p.choice, p.list, hintLabel)
}

func (p *PageFilter) Choice(chosen string) {
	log.Println("PageFilter->Choice: ", chosen)
}

//ParseCIDR: takes a string representing an IP/mask and returns an IP and an IPNet
//IPNet.Contains: c

func (p *PageFilter) AquireData(config *Config) error {
	var ok bool
	config.Filter.Type, ok = MapFilterTypeLabelFromString[p.choice.Selected]
	if !ok {
		fmt.Println("not OK")
	}
	var netsList []*net.IPNet

	config.Filter.Networks = nil
	for _, s := range strings.Split(p.list.Text, "\n") {
		config.Filter.Networks = append(config.Filter.Networks, s)
		comment := strings.Index(s, "#")
		if comment != -1 {
			s = s[:comment]
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		_, net, err := net.ParseCIDR(s)
		if err != nil {
			return err
		}
		netsList = append(netsList, net)
	}
	p.wiz.model = NewModel().SetAccept(func(s string) bool {
		if config.Filter.Type == FilterTypeNoFilter {
			return true
		}
		for _, n := range netsList {
			ip := net.ParseIP(s)
			if ip == nil {
				return false
			}
			if n.Contains(ip) {
				return config.Filter.Type == FilterTypeInclude
			}
		}
		return config.Filter.Type == FilterTypeExclude
	})
	return nil
}
