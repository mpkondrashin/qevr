/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

wizard.go

QeVR wizard framework
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kdeconinck/camelcase"
)

//go:generate fyne bundle --name IconSVGResource --output resource.go image.png

var ErrAbort = errors.New("abort")

type Page interface {
	Index() PageIndex
	Content() fyne.CanvasObject
	Run()
	AquireData(conf *Config) error
	Next(previousPage PageIndex) PageIndex
	Prev() PageIndex
}

func PageName(p Page) string {
	return strings.Join(camelcase.Split(p.Index().String()[2:]), " ")
}

type BasePage struct {
	wiz          *Wizard
	previousPage PageIndex
}

func NewBasePage(wiz *Wizard) BasePage {
	return BasePage{
		wiz: wiz,
	}
}
func (p *BasePage) SavePrevious(previousPage PageIndex) {
	if previousPage == pgExit {
		return
	}
	p.previousPage = previousPage
}

func (p *BasePage) Prev() PageIndex {
	return p.previousPage
}

func (p *BasePage) Run() {}

func (p *BasePage) AquireData(config *Config) error {
	return nil
}

//go:generate stringer -type=PageIndex

type PageIndex int

const (
	pgStart PageIndex = -1
	pgIntro PageIndex = iota
	pgSource
	pgFilter
	pgLoad
	pgTarget
	pgSave
	pgSMS
	pgUpload
	pgFinish
	pgExit
)

type Wizard struct {
	config      *Config
	pages       []Page
	firstPage   PageIndex
	currentPage PageIndex
	app         fyne.App
	win         fyne.Window
	pagesList   *fyne.Container
	buttonsLine *fyne.Container
	csv         *CSV
	model       *Model
}

func NewWizard(config *Config) *Wizard {
	w := &Wizard{
		config:      config,
		app:         app.NewWithID("github.com/mpkondrashin/opsalyzer-release"),
		pagesList:   container.NewVBox(),
		buttonsLine: container.NewHBox(),
	}
	w.app.Lifecycle()
	w.win = w.app.NewWindow("QeVR")
	w.win.Resize(fyne.NewSize(600, 400))
	w.win.SetMaster()
	w.firstPage = w.Pages()
	w.currentPage = w.firstPage
	w.win.SetContent(w.Window())
	return w
}

func (w *Wizard) Pages() PageIndex {
	w.pages = make([]Page, pgExit)
	w.pages[pgIntro] = &PageIntro{BasePage: NewBasePage(w)}
	w.pages[pgSource] = &PageSource{BasePage: NewBasePage(w)}
	w.pages[pgFilter] = &PageFilter{BasePage: NewBasePage(w)}
	w.pages[pgLoad] = &PageLoad{BasePage: NewBasePage(w)}
	w.pages[pgTarget] = &PageTarget{BasePage: NewBasePage(w)}
	w.pages[pgSave] = &PageOutput{BasePage: NewBasePage(w)}
	w.pages[pgSMS] = &PageSMS{BasePage: NewBasePage(w)}
	w.pages[pgUpload] = &PageUpload{BasePage: NewBasePage(w)}
	w.pages[pgFinish] = &PageFinish{BasePage: NewBasePage(w)}
	return pgIntro
}

func (c *Wizard) Window() fyne.CanvasObject {
	if c.currentPage < pgIntro {
		c.currentPage = pgIntro
		log.Printf("Wrong current page: %v (%d)", c.currentPage, c.currentPage)
	}
	if c.currentPage >= pgFinish {
		c.currentPage = pgFinish
		log.Printf("Wrong current page: %v (%d)", c.currentPage, c.currentPage)
	}
	p := c.pages[c.currentPage]
	c.UpdatePagesList()
	middle := container.NewPadded(container.NewVBox(layout.NewSpacer(), p.Content(), layout.NewSpacer()))
	upper := container.NewBorder(nil, nil, container.NewHBox(c.pagesList, widget.NewSeparator()), nil, middle)
	buttons := container.NewBorder(nil, nil, nil, c.buttonsLine)
	bottom := container.NewVBox(widget.NewSeparator(), buttons)
	return container.NewBorder(nil, container.NewPadded(bottom), nil, nil, upper)
}

func (c *Wizard) UpdatePagesList() {
	c.pagesList.RemoveAll()
	image := canvas.NewImageFromResource(Logo)
	image.SetMinSize(fyne.NewSize(52, 52))
	image.FillMode = canvas.ImageFillContain
	c.pagesList.Add(image)
	previous := pgStart
	i := c.firstPage
	for {
		if i == pgExit {
			break
		}
		pg := c.pages[i]
		next := pg.Next(previous)
		if next <= previous {
			panic(fmt.Errorf("pages cycle: %v (%d) -> %v (%d)", previous, previous, next, next))
		}
		if i == c.currentPage {
			c.pagesList.Add(widget.NewLabelWithStyle("▶ "+PageName(pg), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			prev, next := c.Buttons(i == c.firstPage, next == pgExit)
			c.buttonsLine.RemoveAll()
			c.buttonsLine.Add(prev)
			c.buttonsLine.Add(next)
		} else {
			c.pagesList.Add(widget.NewLabel("    " + PageName(pg)))
		}
		previous = i
		i = next
	}
}

func (c *Wizard) Buttons(first, last bool) (*widget.Button, *widget.Button) {
	prevButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), c.Prev)
	if first {
		prevButton.Disable()
	}

	nextButton := widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), c.Next)
	nextButton.IconPlacement = widget.ButtonIconTrailingText

	if last {
		nextButton = widget.NewButtonWithIcon("Quit", theme.CancelIcon(), c.Quit)
	}
	return prevButton, nextButton
}

func (c *Wizard) Quit() {
	log.Print("Quit")
	err := c.pages[c.currentPage].AquireData(c.config)
	if err != nil {
		log.Printf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
	}
	//dialog.ShowConfirm("QeVR", "Exit?", )
	c.app.Quit()
}

func (c *Wizard) Next() {
	log.Printf("Next from page %d", c.currentPage)
	err := c.pages[c.currentPage].AquireData(c.config)
	if err != nil {
		if errors.Is(err, ErrAbort) {
			c.app.Quit()
		}
		log.Printf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
		return
	}
	c.currentPage = c.pages[c.currentPage].Next(c.currentPage)
	c.win.SetContent(c.Window())
	c.pages[c.currentPage].Run()
}

func (c *Wizard) Prev() {
	log.Printf("Prev from page %d to %d", c.currentPage, c.pages[c.currentPage].Prev())
	c.currentPage = c.pages[c.currentPage].Prev()
	c.win.SetContent(c.Window())
	c.pages[c.currentPage].Run()
}

func (c *Wizard) Run() {
	c.win.ShowAndRun()
}
