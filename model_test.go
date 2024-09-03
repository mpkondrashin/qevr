/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

model_test.go

Testing model
*/
package main

import (
	"strings"
	"testing"
)

func TestModelAdd(t *testing.T) {
	m := NewModel()
	m.AddIPAndCVE("1.1.1.1", "CVE-1-1")
	var sb strings.Builder
	m.Save(&sb)
	expected := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-1","MEDIUM"
`
	actual := sb.String()
	if actual != expected {
		t.Errorf("expected \"%s\", but got \"%s\"", expected, actual)
	}
}

func TestModelAddTwice(t *testing.T) {
	m := NewModel()
	m.AddIPAndCVE("1.1.1.1", "CVE-1-1")
	m.AddIPAndCVE("1.1.1.1", "CVE-1-1")
	var sb strings.Builder
	m.Save(&sb)
	expected := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-1","MEDIUM"
`
	actual := sb.String()
	if actual != expected {
		t.Errorf("expected \"%s\", but got \"%s\"", expected, actual)
	}
}

func TestModelAddCombineCVEs(t *testing.T) {
	m := NewModel()
	m.AddIPAndCVE("1.1.1.1", "CVE-1-1")
	m.AddIPAndCVE("1.1.1.1", "CVE-1-2")
	var sb strings.Builder
	m.Save(&sb)
	expected1 := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-1,CVE-1-2","MEDIUM"
`
	expected2 := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-2,CVE-1-1","MEDIUM"
`
	actual := sb.String()
	if actual != expected1 && actual != expected2 {
		t.Errorf("expected \"%s\" or \"%s\", but got \"%s\"", expected1, expected2, actual)
	}
}

func TestModelSplitCVEs(t *testing.T) {
	m := NewModel().SetMaxCVEs(1)
	m.AddIPAndCVE("1.1.1.1", "CVE-1-1")
	m.AddIPAndCVE("1.1.1.1", "CVE-1-2")
	var sb strings.Builder
	m.Save(&sb)
	expected1 := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-1","MEDIUM"
"1.1.1.1","CVE-1-2","MEDIUM"
`
	expected2 := `"IP_ADDRESS","CVE_IDS","SEVERITY"
"1.1.1.1","CVE-1-2","MEDIUM"
"1.1.1.1","CVE-1-1","MEDIUM"
`
	actual := sb.String()
	if actual != expected1 && actual != expected2 {
		t.Errorf("expected \"%s\" or \"%s\", but got \"%s\"", expected1, expected2, actual)
	}
}
