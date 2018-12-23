package config

type WebServerParams struct {
	Production     bool
	Bind           string
	StaticDir      string
	DebugMode      bool
	SSL            bool
	CertCache      string
	HostsWhitelist []string
	Registration   bool       // enable disable registration
	InviteOnly     bool       // allow registration only using valid invite code.
}
