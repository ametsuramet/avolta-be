package config

type ServerConfiguration struct {
	SecretKey    string
	Mode         string
	Addr         string
	Environment  string
	BaseURL      string
	ExpiredToken int
}
