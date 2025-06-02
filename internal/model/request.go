package model

type MatchRequestTemplate struct {
	MustMethod          Method         `yaml:"MustMethod" json:"MustMethod"`
	MustHeaders         map[string]any `yaml:"MustHeaders" json:"MustHeaders"`
	MustPathParameters  map[string]any `yaml:"MustPathParameters" json:"MustPathParameters"`
	MustQueryParameters map[string]any `yaml:"MustQueryParameters" json:"MustQueryParameters"`
	MustBody            map[string]any `yaml:"MustBodyParameters" json:"MustBodyParameters"`
}

type Request struct {
	Path    map[string]any
	Query   map[string]any
	Headers map[string]any
	Body    map[string]any
}
