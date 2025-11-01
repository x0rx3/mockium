package model

type HandleTemplate struct {
	MatchRequestTemplate MatchRequestTemplate `yaml:"MatchRequest" json:"MatchRequest"`
	SetResponseTemplate  SetResponseTemplate  `yaml:"SetResponse" json:"SetResponse"`
}

type Handle struct {
	MatchRequest MatchRequest
	SetResponse  SetResponse
}
