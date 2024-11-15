package database

type Deployment struct {
	Id       string
	Kind     string
	Name     string
	Endpoint string
}

type Navigation struct {
	Id           string
	DeploymentId string
	Endpoint     string
	Title        string
}

type Proxy struct {
	Id             string
	DeploymentId   string
	BackendCode    string
	BackendAddress string
}
