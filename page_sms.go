/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_sms.go

Get SMS access credentials
*/
package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/asaskevich/govalidator"
)

type PageSMS struct {
	BasePage
	addressEntry    *widget.Entry
	ignoreTLSErrors *widget.Check
	apiKeyEntry     *widget.Entry
}

var _ Page = &PageSource{}

func (p *PageSMS) Index() PageIndex {
	return pgSMS
}

func (p *PageSMS) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgUpload
}

func (p *PageSMS) Content() fyne.CanvasObject {
	p.addressEntry = widget.NewEntry()
	p.addressEntry.SetText(p.wiz.config.Output.SMS.Address)
	addressItem := widget.NewFormItem("SMS Server Address:", p.addressEntry)

	p.ignoreTLSErrors = widget.NewCheck("Ignore", nil)
	p.ignoreTLSErrors.Checked = p.wiz.config.Output.SMS.IgnoreTLSErrors
	tlsItem := widget.NewFormItem("On TLS errors:", p.ignoreTLSErrors)

	p.apiKeyEntry = widget.NewEntry()
	p.apiKeyEntry.SetText(p.wiz.config.Output.SMS.APIKey)
	apiKeyItem := widget.NewFormItem("API Key:", p.apiKeyEntry)

	return widget.NewForm(
		addressItem, tlsItem, apiKeyItem,
	)
}

func (p *PageSMS) AquireData(config *Config) error {
	if p.addressEntry.Text == "" {
		return fmt.Errorf("SMS Server Address cannot be empty")
	}
	p.wiz.config.Output.SMS.Address = strings.TrimSpace(p.addressEntry.Text)
	if !govalidator.IsUUID(strings.ToLower(p.apiKeyEntry.Text)) {
		return fmt.Errorf("API Key must be a valid UUID")
	}
	p.wiz.config.Output.SMS.APIKey = strings.TrimSpace(p.apiKeyEntry.Text)

	p.wiz.config.Output.SMS.IgnoreTLSErrors = p.ignoreTLSErrors.Checked
	return nil
	return p.Upload()
}

func (p *PageSMS) Upload() error {
	runTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	url := fmt.Sprintf("https://%s/vulnscanner/import?vendor=SMS-Standard&product=QeVR&version=1&runtime=%s/", p.wiz.config.Output.SMS.Address, runTime)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileName := fmt.Sprintf("qevr_%s.csv", time.Now().Format("20060102"))
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}

	err = p.wiz.model.Save(formFile)
	if err != nil {
		return err
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
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
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("error code: %d\n%s", resp.StatusCode, string(body))
	}
	return nil
}
