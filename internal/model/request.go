package model

import "gomock/internal/transport/method"

type MatchRequestTemplate struct {
	MustMethod          method.Method  `yaml:"MustMethod" json:"MustMethod"`
	MustPathParameters  map[string]any `yaml:"MustPathParameters" json:"MustPathParameters"`
	MustQueryParameters map[string]any `yaml:"MustQueryParameters" json:"MustQueryParameters"`
	MustFormParameters  map[string]any `yaml:"MustFormParameters" json:"MustFormParameters"`
	MustHeaders         map[string]any `yaml:"MustHeaders" json:"MustHeaders"`
	MustBody            map[string]any `yaml:"MustBody" json:"MustBody"`
}

type MatchRequest struct {
	MustMethod          method.Method
	MustPathParameters  map[string]any
	MustQueryParameters map[string]any
	MustFormParameters  map[string]any
	MustHeaders         map[string]any
	MustBody            map[string]any
}
