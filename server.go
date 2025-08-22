package main

import (
	_ "embed"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

func SetDevicePendingAction(ip string, action Action) {
	newDevs := []Device{}
	for _, v := range devices {
		if v.Ip == ip {
			v.PendingAction = action
		}
		newDevs = append(newDevs, v)
	}
	devices = newDevs
}

//from the point of view of the device
var uploads = map[string]*func(name string) io.WriteCloser{}
var downloads = map[string]io.ReadSeekCloser{}

func StartServer() error {
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write(transferPage)
	})

	http.HandleFunc("GET /file/download/{id}", func(w http.ResponseWriter, r *http.Request) {
		reader, ok := downloads[r.PathValue("id")]
		if !ok {
			w.Write([]byte("download failed: internal error"))
			return
		}

		reader.Seek(0, io.SeekStart)
		
		_, err := io.Copy(w, reader)
		if err != nil {
			w.Write([]byte("download failed: could not copy file to a response"))
			return
		}
	})

	http.HandleFunc("POST /file/upload/{id}", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("upload failed: grabbing file: " + err.Error()))
			return
		}
		defer file.Close()

		fptr, ok := uploads[r.PathValue("id")]
		if !ok {
			w.Write([]byte("upload failed: internal error"))
			return
		}
		writer := (*fptr)(header.Filename)
		
		_, cerr := io.Copy(writer, file)
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

				switch v.PendingAction.Type {
				case "ToDevice":
					jsonResponse.Status = "download"

					id := uuid.New().String()
					downloads[id] = v.PendingAction.UploadReader
					jsonResponse.Url = "/file/download/" + id
					jsonResponse.FileName = v.PendingAction.FileName
				case "FromDevice":
					jsonResponse.Status = "upload"
					id := uuid.New().String()
					jsonResponse.Url = "/file/upload/" + id
					uploads[id] = v.PendingAction.DownloadBegins
				}

				v.PendingAction = Action{} //clear out the pending aciton
				foundDevice = true
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