package main

import (
	"io"
	"sync/atomic"
	"time"
)

type Action struct {
	Type           string //either ToDevice or FromDevice. If type is "" the action is ignored
	FileName       string //always set if type is ToDevice, set when the download has begun in case of FromDevice
	UUID string
	Begun bool //Whether the device already has been informed of the download
	BegunTime time.Time //Timestamp when the download begun
	UploadReader   io.ReadSeekCloser //only populated if type is ToDevice
	DownloadBegins *func(name string) io.WriteCloser //only populated if type is FromDevice, called when the download has begun.
	ProgressedWriter *ProgressedWriter //set when the download has begun. Should NOT be written to.
}

type Device struct {
	Ip            string //main identifier
	Name          string //extracted from user agent
	LastPing      int64  //if we've gone more than 3 seconds without a ping, we dont include it
	PendingAction Action
}

type PingResponse struct {
	Status   string `json:"status"`
	FileName string `json:"filename"`
	Url      string `json:"url"`
}
type ProgressedWriter struct {
	writer io.Writer
	closer io.Closer
	progress atomic.Int64
	Size int64 //Size can be set during initalization and may be used to track the expected size of the data to be written to the writer.
}

func GetProgressedWriter(writer io.Writer, size int64) *ProgressedWriter {
	return &ProgressedWriter{
		writer: writer,
		progress: atomic.Int64{},
		Size: size,
	}
}

//BytesWritten returns the amount of bytes written to the ProgressedWriter
func (w *ProgressedWriter) BytesWritten() int64 {
	return w.progress.Load() 
}

//Write calls write on the underlying io.Writer
func (w *ProgressedWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	w.progress.Add(int64(n))
	return
}

//Close calls close on the underlying io.Closer
func (w *ProgressedWriter) Close() error {
	return w.closer.Close()
}