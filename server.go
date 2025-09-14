package main

import (
	_ "embed"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed transfer.html
var transferPage []byte
var devices = []Device{}

func GetValidDevices() []Device {
	curTime := time.Now().Unix()
	validDevices := []Device{}

	for _, v := range devices {
		v.PendingAction = Action{}
		if v.LastPing + 3 > curTime {
			validDevices = append(validDevices, v)
		}
	}
	return validDevices
}

//either provide ip or uuid
func SetDevicePendingAction(ip string, action Action, uuid string) {
	newDevs := []Device{}
	for _, v := range devices {
		if v.Ip == ip || v.PendingAction.UUID == uuid {
			v.PendingAction = action
		}
		newDevs = append(newDevs, v)
	}
	devices = newDevs
}

//either provide ip or uuid
func GetDevicePendingAction(ip string, uuid string) Action {
	for _, v := range devices {
		if v.Ip == ip || v.PendingAction.UUID == uuid {
			return v.PendingAction
		}
	}
	return Action{}
}


//from the point of view of the device
var uploads = map[string]*func(name string) io.WriteCloser{}
var downloads = map[string]io.ReadSeekCloser{}

func StartServer() error {
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write(transferPage)
	})

	http.HandleFunc("GET /file/download/{id}", func(w http.ResponseWriter, r *http.Request) {
		defer SetDevicePendingAction("", Action{}, r.PathValue("id"))

		reader, ok := downloads[r.PathValue("id")]
		if !ok {
			w.Write([]byte("download failed: internal error"))
			return
		}

		end, serr := reader.Seek(0, io.SeekEnd)
		if serr != nil {
			w.Write([]byte("download failed: could not get file length"))
			return
		}
		w.Header().Add("Content-Length", strconv.FormatInt(end, 10))

		reader.Seek(0, io.SeekStart)
		pw := GetProgressedWriter(w, end)
		pw.closer = reader

		na := GetDevicePendingAction("", r.PathValue("id"))
		na.ProgressedWriter = pw
		na.BegunTime = time.Now()
		SetDevicePendingAction("", na, r.PathValue("id"))
		
		_, err := io.Copy(pw, reader)
		if err != nil {
			w.Write([]byte("download failed: could not copy file to a response"))
			return
		}
	})

	http.HandleFunc("POST /file/upload/{id}", func(w http.ResponseWriter, r *http.Request) {
		defer SetDevicePendingAction("", Action{}, r.PathValue("id"))
		na := GetDevicePendingAction("", r.PathValue("id"))
		na.UploadPreWriterStatus = "RequestReceived"
		SetDevicePendingAction("", na, r.PathValue("id"))

		file, header, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("upload failed: grabbing file: " + err.Error()))
			return
		}
		defer file.Close()

		if GetDevicePendingAction("", r.PathValue("id")).UploadPreWriterStatus == "Canceled" {
			w.Write([]byte("upload failed: canceled"))
			return
		}

		fptr, ok := uploads[r.PathValue("id")]
		if !ok {
			w.Write([]byte("upload failed: internal error"))
			return
		}
		writer := (*fptr)(header.Filename)
		defer writer.Close()
		pwriter := GetProgressedWriter(writer, header.Size)
		pwriter.closer = writer

		
		na.ProgressedWriter = pwriter
		na.FileName = header.Filename
		na.BegunTime = time.Now()
		SetDevicePendingAction("", na, r.PathValue("id"))

		_, cerr := io.Copy(pwriter, file)
		if cerr != nil {
			w.Write([]byte("upload failed: could not copy upload to a file"))
			return
		}
		
		w.Header().Add("Location", "/")
		w.WriteHeader(301)
	})

	http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		var jsonResponse = PingResponse{
			Status: "yield",
		}
		var newDevices = []Device{}
		foundDevice := false

		for _, v := range devices {
			if strings.Split(r.RemoteAddr, ":")[0] == v.Ip {
				v.LastPing = time.Now().Unix()
				foundDevice = true

				if !v.PendingAction.Begun {
					switch v.PendingAction.Type {
					case "ToDevice":
						jsonResponse.Status = "download"

						downloads[v.PendingAction.UUID] = v.PendingAction.UploadReader
						jsonResponse.Url = "/file/download/" + v.PendingAction.UUID
						jsonResponse.FileName = v.PendingAction.FileName
					case "FromDevice":
						jsonResponse.Status = "upload"
						jsonResponse.Url = "/file/upload/" + v.PendingAction.UUID
						uploads[v.PendingAction.UUID] = v.PendingAction.DownloadBegins
					}
				}
				v.PendingAction.Begun = true
			}
			newDevices = append(newDevices, v)
		}

		if !foundDevice {
			ua := r.Header.Get("User-Agent")
			ua = strings.Split(ua, "(")[1]
			ua = strings.Split(ua, ")")[0]
			ua = "(" + ua + ")"

			newDevices = append(newDevices, Device{
				Ip: strings.Split(r.RemoteAddr, ":")[0],
				Name: html.EscapeString(strings.Split(r.RemoteAddr, ":")[0] + " " + ua), //192.168.69.420 (Device type)
				LastPing: time.Now().Unix(),
			})
		}
		devices = newDevices
		ba, err := json.Marshal(jsonResponse)
		if err != nil {
			w.Write([]byte("jsonerr"))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(ba)
	})
	var err error
	go func() {
		err = http.ListenAndServe(":5678", http.DefaultServeMux)
	}()
	time.Sleep(200 * time.Millisecond)
	if err != nil {return err}

	return nil
}