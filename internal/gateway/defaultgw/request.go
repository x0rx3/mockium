package defaultgw

import (
	"context"
	"mockium/pkg/core"
)

func NewRequest(ctx context.Context, protocol string, payload []byte, metadata map[string][]string) core.Request {
	return &Request{
		protocol: protocol,
		ctx:      ctx,
		payload:  payload,
		metadata: metadata,
	}
}

type Request struct {
	protocol string
	ctx      context.Context
	payload  []byte
	metadata map[string][]string
}

func (inst *Request) Protocol() string              { return inst.protocol }
func (inst *Request) Context() context.Context      { return inst.ctx }
func (inst *Request) Payload() []byte               { return inst.payload }
func (inst *Request) Metadata() map[string][]string { return inst.metadata }
func (inst *Request) MetadataValue(key string) []string {
	if values, exists := inst.metadata[key]; exists {
		return values
	}

	return nil
}
