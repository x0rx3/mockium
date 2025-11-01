// Package model
package model

type CookieTemplate struct {
	Name     string `json:"Name"`
	Value    string `json:"Value"`
	Path     string `json:"Path"`
	Domain   string `json:"Domain"`
	Expires  string `json:"Expires"`
	MaxAge   string `json:"MaxAge"`
	Secure   bool   `json:"Secure"`
	HttpOnly bool   `json:"HttpOnly"`
	SameSite string `json:"SameSite"`
}

type Cookie struct {
	Name     string
	Value    any
	Path     any
	Domain   any
	Expires  any
	MaxAge   any
	Secure   bool
	HttpOnly bool
	SameSite any
}
