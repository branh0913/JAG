package jobbuilder


type Config interface {
	New(endpoint, apitoken, jconfigpath string) (*JBuilderConfig, error)
	BuildFile(currentuser string)
}


type JBuilderConfig struct {
	APIToken string
	Endpoint string
}


