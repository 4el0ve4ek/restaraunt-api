package jwt

type Config struct {
	SecretKey string `yaml:"secret_key" json:"secret_key"`
}
