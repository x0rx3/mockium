package model

import "time"

type ProcessLoggingFileds struct {
	Time     time.Time      `json:"time"`
	Request  *LogginRequest `json:"request"`
	Response SetResponse    `json:"response"`
}

type LogginRequest struct {
	Url        string         `json:"url"`
	Method     string         `json:"method"`
	RemoteAddr string         `json:"reqmote_addr"`
	Headers    map[string]any `json:"headers"`
	Body       any            `json:"body"`
}
