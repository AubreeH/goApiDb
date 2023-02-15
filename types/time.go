package dataTypes

import (
	"time"
)

type Time struct {
	time.Time
}

func (value Time) ExtractDataFunc() any {
	return value.Format("2006-01-02 15:04:05")
}

func (_ Time) GetPtrFunc(value *Time) *time.Time {
	return &value.Time
}
