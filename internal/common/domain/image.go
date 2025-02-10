package domain

type Image struct {
	URL    string `json:"url,omitempty"`
	Base64 string `json:"base64,omitempty"`
}
