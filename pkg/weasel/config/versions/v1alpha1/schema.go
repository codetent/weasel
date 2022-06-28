package v1alpha1

const Version string = "v1alpha1"

type Config struct {
	Version      string
	Name         string
	Environments map[string]*Environment
}

type Environment struct {
	Image string
}
