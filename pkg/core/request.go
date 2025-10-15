package core

import "context"

type Request interface {
	Protocol() string
	Context() context.Context
	Payload() []byte
	Metadata() map[string][]string
	MetadataValue(string) []string
}
