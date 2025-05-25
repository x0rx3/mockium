package model

type Template struct {
	Path   string           `yaml:"Path" json:"Path"`
	Handle []HandleTemplate `yaml:"Handle" json:"Handle"`
}
