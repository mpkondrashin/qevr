/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_upload.go

Upload eVR CSV to SMS
*/
package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageUpload struct {
	BasePage
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
}

var _ Page = &PageUpload{}

func (p *PageUpload) Index() PageIndex {
	return pgUpload
}

func (p *PageUpload) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgFinish
}

func (p *PageUpload) Content() fyne.CanvasObject {
	p.progressBar = widget.NewProgressBar()
	p.statusLabel = widget.NewLabel("")
	return container.NewVBox(
		p.progressBar,
		p.statusLabel,
	)
}

func (p *PageUpload) Run() {
	p.statusLabel.SetText("Uploading...")
	quit := p.Progress()
	fileName, err := p.Upload()
	quit()
	if err != nil {
		log.Print(err)
		p.statusLabel.SetText("Error")
		dialog.ShowError(err, p.wiz.win)
	}
	p.progressBar.SetValue(1.0)
	p.statusLabel.SetText(fmt.Sprintf("Vulnerability scan report successfully uploaded as %s\nGo to SMS console -> Profiles-> Vulnerability Scans (eVR)", fileName))
}

func (p *PageUpload) Progress() func() {
	quit := make(chan struct{})
	go func() {
		const steps = 100
		for i := 1; i < steps; i += 1 {
			select {
			case <-quit:
				return
			default:
				time.Sleep(p.wiz.config.Output.SMS.Timeout / steps)
				p.progressBar.SetValue(float64(i) / float64(steps))
			}
		}
	}()
	return func() {
		quit <- struct{}{}
	}
}

func (p *PageUpload) Upload() (string, error) {
	runTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	url := fmt.Sprintf("https://%s/vulnscanner/import?vendor=SMS-Standard&product=QeVR&version=1&runtime=%s/", p.wiz.config.Output.SMS.Address, runTime)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileName := fmt.Sprintf("qevr_%s.csv", time.Now().Format("20060102"))
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	err = p.wiz.model.Save(formFile)
	if err != nil {
		return "", err
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-SMS-API-KEY", p.wiz.config.Output.SMS.APIKey)
	// Ignore TLS errors
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: p.wiz.config.Output.SMS.IgnoreTLSErrors},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   p.wiz.config.Output.SMS.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		return "", fmt.Errorf("error code: %d\n%s", resp.StatusCode, string(body))
	}
	return fileName, nil
}

func (p *PageUpload) AquireData(config *Config) error {
	return nil
}
