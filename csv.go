/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

csv.go

Input CSV file management
*/
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
)

type CSV struct {
	fileName string
	IPIndex  int
	CVEIndex int
	Header   bool
}

func NewCSV(fileName string) *CSV {
	return &CSV{
		fileName: fileName,
	}
}

func (c *CSV) Load(progress func(int), callback func(ip, cve string)) error {
	if c.IPIndex == 0 && c.CVEIndex == 0 {
		return fmt.Errorf("run DetectIPCVE first")
	}
	file, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	progressReader := NewProgressReader(file, progress)
	reader := csv.NewReader(progressReader)
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}
		if SkipHeader && lineCount == 0 {
			lineCount++
			continue
		}
		if len(record) > c.IPIndex && len(record) > c.CVEIndex {
			callback(record[c.IPIndex], record[c.CVEIndex])
		}
		lineCount++
	}
	return nil
}

const (
	IPAddress   = `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?(\.|$)){4}$`
	CVEID       = `^CVE-\d{4}-\d{4,7}$`
	SampleLines = 10
	SkipHeader  = true
)

func (c *CSV) DetectIPCVE() (err error) {
	lines, err := c.SampleCSV(SampleLines)
	if err != nil {
		return err
	}
	c.IPIndex, c.CVEIndex, err = Detect(lines, IPAddress, CVEID)
	return
}

func (c *CSV) SampleCSV(sampleLines int) ([][]string, error) {
	file, err := os.Open(c.fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var lines [][]string
	for i := 0; i < sampleLines; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Detect CSV structure: %w", err)
		}
		if i == 0 && SkipHeader {
			continue
		}
		lines = append(lines, record)
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data found in the file")
	}
	return lines, nil
}

func Detect(lines [][]string, m1, m2 string) (int, int, error) {
	r1 := regexp.MustCompile(m1)
	r2 := regexp.MustCompile(m2)
	p1 := make(map[int]int)
	p2 := make(map[int]int)
	for _, line := range lines {
		for i := 0; i < len(line); i++ {
			if r1.Match([]byte(line[i])) {
				p1[i]++
			}
			if r2.Match([]byte(line[i])) {
				p2[i]++
			}
		}
	}
	c1 := Max(p1)
	if c1 == -1 {
		return 0, 0, fmt.Errorf("\"%s\" not found", m1)
	}
	c2 := Max(p2)
	if c2 == -1 {
		return 0, 0, fmt.Errorf("\"%s\" not found", m2)
	}
	return c1, c2, nil
}

func Max(m map[int]int) int {
	maxCount := 0
	maxIndex := -1
	for index, count := range m {
		if count > maxCount {
			maxCount = count
			maxIndex = index
		}
	}
	return maxIndex
}
