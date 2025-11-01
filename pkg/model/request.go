package model

type MatchRequestTemplate struct {
	MustMethod          string              `yaml:"MustMethod" json:"MustMethod"`
	MustType            string              `yaml:"MustType" json:"MustType"` // unary, stream
	MustHeaders         *MustHeaderTemplate `yaml:"MustHeaders" json:"MustHeaders"`
	MustPathParameters  map[string]any      `yaml:"MustPathParameters" json:"MustPathParameters"`
	MustQueryParameters map[string]any      `yaml:"MustQueryParameters" json:"MustQueryParameters"`
	MustBody            map[string]any      `yaml:"MustBodyParameters" json:"MustBodyParameters"`
}

type MatchRequest struct {
	MustType            string // unary, stream
	MustMethod          string
	MustHeaders         *MustHeader
	MustPathParameters  map[string]any
	MustQueryParameters map[string]any
	MustBody            map[string]any
}
