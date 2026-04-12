package model

type APIErrorResponse struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Path    string `json:"path,omitempty"`
	Details string `json:"details,omitempty"`
}
