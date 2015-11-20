package models

const (
	STATUS_READY = iota
	STATUS_DOWNLOADING
	STATUS_FAIL
	STATUS_FINISH
)

type STATUS int

var GLOBAL_TASK_ID_INCREMENTOR = 1

type Task struct {
	Id int
	Url string
	Status STATUS
	DownloadLink string // mapping to the file system
}

