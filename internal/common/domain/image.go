package domain

type Image struct {
	Url    string `json:"url,omitempty"`
	Base64 string `json:"base64,omitempty"`
}
