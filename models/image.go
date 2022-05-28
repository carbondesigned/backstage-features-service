package models

import (
	"os"
)

type Image struct {
	UserId string
	File   *os.File
	Url    string
}
