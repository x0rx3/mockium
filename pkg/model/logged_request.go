package model

type LoggedRequest struct {
	URL        string         `json:"url"`
	Method     string         `json:"method"`
	RemoteAddr string         `json:"reqmote_addr"`
	Headers    map[string]any `json:"headers"`
	Body       any            `json:"body"`
}
