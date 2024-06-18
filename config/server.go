package config

type ServerConfiguration struct {
	SecretKey    string
	Mode         string
	Addr         string
	Environment  string
	BaseURL      string
	FrontendURL  string
	CmsURL       string
	ExpiredToken int
}
