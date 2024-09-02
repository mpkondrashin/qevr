/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

main.go

Installer main file
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/panicwrap"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func FatalWarning(title, text string) {
	//dialog.ShowInformation("Information", "This is a sample message", s.window)
	a := app.New()
	w := a.NewWindow("QeVR Error: " + title)
	//w.SetMaster() // will it exit the application?
	w.SetContent(container.NewVBox(
		widget.NewLabel(text),
		widget.NewButton("Ok", func() { a.Quit() }),
	))
	w.ShowAndRun()
	os.Exit(1)
}

const qevrLog = "qevr.log"

func SetupLogging(logFileName string) (func(), error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}
	logFolder := filepath.Dir(path)
	logFilename := filepath.Join(logFolder, logFileName)
	file, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.Println("Logging Started")

	exitStatus, err := panicwrap.BasicWrap(func(output string) {
		log.Println(output)
		os.Exit(1)
	})
	if err != nil {
		panic(err)
	}
	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}
	return func() {
		log.Print("Close Logging")
		file.Close()
	}, nil
}

func main() {
	close, err := SetupLogging(qevrLog)
	if err != nil {
		msg := fmt.Sprintf("NewFileLog: %v", err)
		fmt.Fprintln(os.Stderr, msg)
		FatalWarning("SetupLogging Error", msg)
	}
	defer close()
	//log.Printf("Start. Version %s Build %s", globals.Version, globals.Build)
	log.Printf("OS: %s (%s)", runtime.GOOS, runtime.GOARCH)
	log.Print("Starting Wizard")
	configPath, err := ConfigFilePath()
	if err != nil {
		FatalWarning("load config", err.Error())
	}
	config, err := LoadConfig(configPath)
	if err != nil {
		FatalWarning("Load config", err.Error())
	}
	c := NewWizard(config)
	c.Run()
	log.Print("Setup finished")
}

func ExecutableFolder() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	execFolder := filepath.Dir(execPath)
	return execFolder, nil
}
