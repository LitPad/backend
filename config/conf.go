package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ProjectName               string `mapstructure:"PROJECT_NAME"`
	Debug                     bool   `mapstructure:"DEBUG"`
	EmailOtpExpireSeconds     int64  `mapstructure:"EMAIL_OTP_EXPIRE_SECONDS"`
	AccessTokenExpireMinutes  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpireMinutes int    `mapstructure:"REFRESH_TOKEN_EXPIRE_MINUTES"`
	Port                      string `mapstructure:"PORT"`
	SecretKey                 string `mapstructure:"SECRET_KEY"`
	FirstSuperuserEmail       string `mapstructure:"FIRST_SUPERUSER_EMAIL"`
	FirstSuperUserPassword    string `mapstructure:"FIRST_SUPERUSER_PASSWORD"`
	FirstWriterEmail          string `mapstructure:"FIRST_WRITER_EMAIL"`
	FirstWriterPassword       string `mapstructure:"FIRST_WRITER_PASSWORD"`
	FirstReaderEmail          string `mapstructure:"FIRST_READER_EMAIL"`
	FirstReaderPassword       string `mapstructure:"FIRST_READER_PASSWORD"`
	PostgresUser              string `mapstructure:"POSTGRES_USER"`
	PostgresPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresServer            string `mapstructure:"POSTGRES_SERVER"`
	PostgresPort              string `mapstructure:"POSTGRES_PORT"`
	PostgresDB                string `mapstructure:"POSTGRES_DB"`
	TestPostgresDB            string `mapstructure:"TEST_POSTGRES_DB"`
	MailSenderEmail           string `mapstructure:"MAIL_SENDER_EMAIL"`
	MailSenderPassword        string `mapstructure:"MAIL_SENDER_PASSWORD"`
	MailSenderHost            string `mapstructure:"MAIL_SENDER_HOST"`
	MailSenderPort            int    `mapstructure:"MAIL_SENDER_PORT"`
	CORSAllowedOrigins        string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	GoogleClientID            string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret        string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	FacebookAppID             string `mapstructure:"FACEBOOK_APP_ID"`
	SocialsPassword           string `mapstructure:"SOCIALS_PASSWORD"`
}

func GetConfig() (config Config) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	var err error
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
	return
}
