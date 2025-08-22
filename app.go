package main

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetValidDevices() []Device {
	return GetValidDevices()
}

func (a *App) SendFile(deviceIp string) {
	dialog := application.OpenFileDialog()
	dialog.SetTitle("Choose a file to transfer")
	dialog.CanChooseDirectories(false)

	fpath, err := dialog.PromptForSingleSelection()
	if err != nil || fpath == "" {return}

	file, err := os.Open(fpath)
	if err != nil {return}

	SetDevicePendingAction(deviceIp, Action{
		Type: "ToDevice",
		FileName: filepath.Base(fpath),
		UploadReader: file,
	})
}

func (a *App) RecvFile(deviceIp string) {
	dcb := func (name string) io.WriteCloser {
		dialog := application.SaveFileDialog()
		dialog.SetMessage("Choose a place to save the received file.")
		dialog.SetFilename(name)

		fpath, err := dialog.PromptForSingleSelection()
		if err != nil {return nil}

		file, err := os.Create(fpath)
		if err != nil {return nil}

		return file
	}
	SetDevicePendingAction(deviceIp, Action{
		Type: "FromDevice",
		DownloadBegins: &dcb,
	})
}

func  (a *App) GetLocalIp() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "IP"
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP.String()
}