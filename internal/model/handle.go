package model

type HandleTemplate struct {
	MatchRequestTemplate MatchRequestTemplate `yaml:"MatchRequest" json:"MatchRequest"`
	SetResponseTemplate  SetResponseTemplate  `yaml:"SetResponse" json:"SetResponse"`
}
