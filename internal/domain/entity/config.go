package entity

type ConfigService struct {
	Server     ServerConfig     `json:"server" binding:"required"`
	BusinessDB BusinessDBConfig `envPrefix:"POSTGRES_"`
	Ai         AiConfig         `envPrefix:"AI_"`
	Minio      MinioConfig      `envPrefix:"MINIO_"`
}

type ServerConfig struct {
	Port       string `json:"server_port" binding:"required"`
	ServerMode string `json:"server_mode" binding:"required"`
	Domain     string `json:"server_domain" binding:"required"`
}

type BusinessDBConfig struct {
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USER"`
	DBName   string `env:"DB"`
	SSLMode  string `env:"SSL_MODE"`
}

type MinioConfig struct {
	Endpoint        string `env:"ENDPOINT"`
	User            string `env:"ROOT_USER"`
	Password        string `env:"ROOT_PASSWORD"`
	SSL             bool   `env:"USE_SSL"`
	MinioBucketName string `env:"BUCKET_NAME"`
}
type AiConfig struct {
	Model  string `env:"MODEL"`
	ApiKey string `env:"API_KEY"`
}
