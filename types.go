package main

type Action struct {
	Type string //either ToDevice or FromDevice. If type is "" the action is ignored
	FileName string //only populated if type is ToDevice
	DownloadCallback *func(name string, data []byte) //only populated if type is FromDevice
	Data []byte //only populated if type is ToDevice
}

type Device struct {
	Ip string //main identifier
	Name string //extracted from user agent
	LastPing int64 //if we've gone more than 3 seconds without a ping, we dont include it
	PendingAction Action
}

type PingResponse struct {
	Status string `json:"status"`
	FileName string `json:"filename"`
	Url string `json:"url"`
}