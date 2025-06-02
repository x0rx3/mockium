package service

import "mockium/internal/model"

type Comparer interface {
	Compare(expected, actual any) bool
}

type ProcessLogger interface {
	Log(logReq *model.ProcessLoggingFileds)
}
