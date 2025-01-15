package domain

import "time"

type File struct {
	Id         int64
	FileName   string
	Link       string
	Uploader   string
	UploaderId int64
	CreateAt   time.Time
}
