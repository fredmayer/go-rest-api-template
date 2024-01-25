package models

import "time"

type Dummy struct {
	Id          int32
	Name        string
	Description string
	CreatedAt   time.Time
}
