package config

type WebServerParams struct {
	Bind           string
	StaticDir      string
	DebugMode      bool
	SSL            bool
	CertCache      string
	HostsWhitelist []string
}
