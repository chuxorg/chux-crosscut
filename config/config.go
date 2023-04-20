package config

type IConfig interface {
	LoadConfig() (IConfig, error)
	GetSecret(name string) (string, error)
	GetSecrets() (map[string]string, error)
}
