package model

import "time"

type LogEntry struct {
	Time     time.Time      `json:"time"`
	Request  *LoggedRequest `json:"request"`
	Response SetResponse    `json:"response"`
}
