package main

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// App struct
type App struct {
	app *application.App
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) setApp(app *application.App) {
	a.app = app
}

func (a *App) GetValidDevices() []Device {
	return GetValidDevices()
}

func (a *App) openTransferWindow(uuid string) {
	w := a.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "XPSync",
		DisableResize: true,
		Width: 378,
		Height: 250,
		Frameless: true,
		DevToolsEnabled: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(236, 233, 216),
		URL:              "/transfer.html",
	})
	w.OnWindowEvent(events.Common.WindowRuntimeReady, func(event *application.WindowEvent) {
		w.ExecJS("let transferUUID = '" + uuid + "';")
	})
}

func (a *App) CloseTransfer(uuid string) {
	for _, dev := range devices {
		if dev.PendingAction.UUID == uuid && dev.PendingAction.ProgressedWriter != nil {
			dev.PendingAction.ProgressedWriter.Close()
		}
		dev.PendingAction.UploadPreWriterStatus = "Canceled"
	}
}

func (a *App) GetTransferStatus(uuid string) (bytesCopied int64, totalBytes int64, progress float64, message string, fileName string, device string, complete bool, bytesPerSecond float64) {
	bytesCopied = 0
	totalBytes = 0
	progress = float64(1)
	message = "Transfer complete"
	fileName = ""
	device = ""
	bytesPerSecond = 0
	complete = true

	for _, dev := range devices {
		if dev.PendingAction.UUID == uuid {
			pw := dev.PendingAction.ProgressedWriter
			device = dev.Ip
			if pw == nil && dev.PendingAction.UploadPreWriterStatus == "RequestReceived" {
				message = "Preparing"
				progress = -1
				fileName = "<Unknown>"
			} else if pw == nil {
				message = "Waiting for device"
				progress = 0
				fileName = "<Unknown>"
			} else {
				message = "Transferring"
				fileName = dev.PendingAction.FileName
				totalBytes = pw.Size
				bytesCopied = pw.BytesWritten()
				progress = float64(bytesCopied) / float64(totalBytes)
				bytesPerSecond = float64(bytesCopied) / time.Since(dev.PendingAction.BegunTime).Seconds()
			}
			complete = false
		}
	}
	return
}

func (a *App) SendFile(deviceIp string) {
	dialog := application.OpenFileDialog()
	dialog.SetTitle("Choose a file to transfer")
	dialog.CanChooseDirectories(false)

	fpath, err := dialog.PromptForSingleSelection()
	if err != nil || fpath == "" {return}

	file, err := os.Open(fpath)
	if err != nil {return}

	uuid := uuid.New().String()
	SetDevicePendingAction(deviceIp, Action{
		Type: "ToDevice",
		FileName: filepath.Base(fpath),
		UploadReader: file,
		UUID: uuid,
	}, "")
	a.openTransferWindow(uuid)
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

	uuid := uuid.New().String()
	SetDevicePendingAction(deviceIp, Action{
		Type: "FromDevice",
		DownloadBegins: &dcb,
		UUID: uuid,
	}, "")
	a.openTransferWindow(uuid)
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