package config

type Config struct {
	DBUsername   string `mapstructure:"DB_USERNAME"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBName       string `mapstructure:"DB_NAME"`
	Port         string `mapstructure:"PORT"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	JWTIssuer    string `mapstructure:"JWT_ISSUER"`
	LLMProvider  string `mapstructure:"LLM_PROVIDER"`
	LLMBaseURL   string `mapstructure:"LLM_BASE_URL"`
	LLMAPIKey    string `mapstructure:"LLM_API_KEY"`
	LLMModel     string `mapstructure:"LLM_MODEL"`
	EnableDBLog  bool   `mapstructure:"ENABLE_DB_LOG"`
	Verbose      bool   `mapstructure:"VERBOSE"`
}
