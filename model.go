package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const DefaultMaxCVEsPerIP = 2800

type Accept func(string) bool

type Model struct {
	accept       Accept
	maxCVEsPerIP int
	data         map[string]map[string]struct{}
}

func NewModel() *Model {
	return &Model{
		data:         make(map[string]map[string]struct{}),
		maxCVEsPerIP: DefaultMaxCVEsPerIP,
	}
}

func (m *Model) SetAccept(accept Accept) *Model {
	m.accept = accept
	return m
}

func (m *Model) SetMaxCVEs(maxCVEs int) *Model {
	m.maxCVEsPerIP = maxCVEs
	return m
}

func (m *Model) AddIPAndCVE(ip, cve string) {
	if m.accept != nil && !m.accept(ip) {
		return
	}
	if _, exists := m.data[ip]; !exists {
		m.data[ip] = make(map[string]struct{})
	}
	m.data[ip][cve] = struct{}{}
}

func (m *Model) SaveToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	return m.Save(file)
}

func (m *Model) Save(file io.Writer) error {
	header := "IP_ADDRESS,CVE_IDS,SEVERITY\n"
	_, err := file.Write([]byte(header))
	if err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}
	severity := "MEDIUM"
	for ip, cves := range m.data {
		count := 1
		var cvesList []string
		for cve := range cves {
			cvesList = append(cvesList, cve)
			if count%m.maxCVEsPerIP == 0 || count == len(cves) {
				cvesStr := strings.Join(cvesList, ",")
				record := fmt.Sprintf("\"%s\",\"%s\",\"%s\"\n", ip, cvesStr, severity)
				_, err := file.Write([]byte(record))
				if err != nil {
					return fmt.Errorf("error writing CSV record: %w", err)
				}
				cvesList = nil // Clear the list for the next batch
			}
			count++
		}
	}
	return nil
}

func (m *Model) Status() string {
	return fmt.Sprintf("Loaded %d addresses", len(m.data))
}

func (m *Model) FinalStatus() string {
	totalCVEs := 0
	for _, cves := range m.data {
		totalCVEs += len(cves)
	}
	return fmt.Sprintf("Loaded %d CVEs for %d addresses", totalCVEs, len(m.data))
}
