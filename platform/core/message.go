package core

type DeploymentRequest struct {
	Version     string       `yaml:"Version,omitempty"`
	Kind        string       `yaml:"Kind,omitempty"`
	Endpoint    string       `yaml:"Endpoint,omitempty"`
	Proxies     []Proxy      `yaml:"Proxies,omitempty"`
	Navigations []Navigation `yaml:"Navigations,omitempty"`
}

type Navigation struct {
	Endpoint string `yaml:",omitempty"`
	Title    string `yaml:",omitempty"`
}

type Proxy struct {
	BackendCode    string `yaml:"BackendCode,omitempty"`
	BackendAddress string `yaml:"BackendAddress,omitempty"`
	Secure         bool   `yaml:"Secure,omitempty"`
}
