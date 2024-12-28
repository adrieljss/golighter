package platform

import (
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

type Env struct {
	StageStatus          string `mapstructure:"STAGE_STATUS"`
	ServerAddress        string `mapstructure:"SERVER_ADDRESS"`
	ServerRequestTimeout int    `mapstructure:"SERVER_REQUEST_TIMEOUT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	JWTSecretAccessToken  string        `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	JWTSecretRefreshToken string        `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
	JWTAccessTokenTTL     time.Duration `mapstructure:"JWT_ACCESS_TOKEN_TTL"`
	JWTRefreshTokenTTL    time.Duration `mapstructure:"JWT_REFRESH_TOKEN_TTL"`
}

func NewEnv(configFile string) *Env {
	env := &Env{}

	// recursively search for the .env file
	// use this because in testing mode, the
	// working directory is the test directory
	// and not the root directory
	maxDepth := 10
	for i := 0; i < maxDepth; i++ {
		viper.SetConfigFile(configFile)
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err == nil {
			break
		}

		configFile = "../" + configFile
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot find the .env file: ", err)
	}

	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal("error unmarshalling environment variables: ", err)
	}

	if env.StageStatus != "dev" && env.StageStatus != "prod" {
		log.Fatal("invalid app stage: ", env.StageStatus)
	}

	log.Info("current app stage set to: ", env.StageStatus)

	return env
}
