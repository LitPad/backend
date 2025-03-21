package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ProjectName               string `mapstructure:"PROJECT_NAME"`
	Environment               string `mapstructure:"ENVIRONMENT"`
	EmailOtpExpireSeconds     int64  `mapstructure:"EMAIL_OTP_EXPIRE_SECONDS"`
	AccessTokenExpireMinutes  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpireMinutes int    `mapstructure:"REFRESH_TOKEN_EXPIRE_MINUTES"`
	Port                      string `mapstructure:"PORT"`
	SecretKey                 string `mapstructure:"SECRET_KEY"`
	SecretKeyByte             []byte
	FirstSuperuserEmail       string `mapstructure:"FIRST_SUPERUSER_EMAIL"`
	FirstSuperUserPassword    string `mapstructure:"FIRST_SUPERUSER_PASSWORD"`
	FirstAuthorEmail          string `mapstructure:"FIRST_AUTHOR_EMAIL"`
	FirstAuthorPassword       string `mapstructure:"FIRST_AUTHOR_PASSWORD"`
	FirstReaderEmail          string `mapstructure:"FIRST_READER_EMAIL"`
	FirstReaderPassword       string `mapstructure:"FIRST_READER_PASSWORD"`
	PostgresUser              string `mapstructure:"POSTGRES_USER"`
	PostgresPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresServer            string `mapstructure:"POSTGRES_SERVER"`
	PostgresPort              string `mapstructure:"POSTGRES_PORT"`
	PostgresDB                string `mapstructure:"POSTGRES_DB"`
	TestPostgresDB            string `mapstructure:"TEST_POSTGRES_DB"`
	MailSenderEmail           string `mapstructure:"MAIL_SENDER_EMAIL"`
	MailFrom                  string `mapstructure:"MAIL_FROM"`
	MailSenderPassword        string `mapstructure:"MAIL_SENDER_PASSWORD"`
	MailSenderHost            string `mapstructure:"MAIL_SENDER_HOST"`
	MailSenderPort            int    `mapstructure:"MAIL_SENDER_PORT"`
	MailApiKey                string `mapstructure:"MAIL_API_KEY"`
	BrevoListID               int    `mapstructure:"BREVO_LIST_ID"`
	BrevoContactsUrl          string `mapstructure:"BREVO_CONTACTS_URL"`
	CORSAllowedOrigins        string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowCredentials      bool   `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	GoogleAndroidClientID     string `mapstructure:"GOOGLE_ANDROID_CLIENT_ID"`
	GoogleIOSClientID         string `mapstructure:"GOOGLE_IOS_CLIENT_ID"`
	FacebookAppID             string `mapstructure:"FACEBOOK_APP_ID"`
	SocialsPassword           string `mapstructure:"SOCIALS_PASSWORD"`
	EmailVerificationPath     string `mapstructure:"EMAIL_VERIFICATION_PATH"`
	PasswordResetPath         string `mapstructure:"PASSWORD_RESET_PATH"`
	StripePublicKey           string `mapstructure:"STRIPE_PUBLIC_KEY"`
	StripeSecretKey           string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret       string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	SocketSecret              string `mapstructure:"SOCKET_SECRET"`
	S3AccessKey               string `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey               string `mapstructure:"S3_SECRET_KEY"`
	S3EndpointUrl             string `mapstructure:"S3_ENDPOINT_URL"`
	BookCoverImagesBucket     string `mapstructure:"BOOK_COVER_IMAGES_BUCKET"`
	UserImagesBucket          string `mapstructure:"USER_IMAGES_BUCKET"`
	IDFrontImagesBucket       string `mapstructure:"ID_FRONT_IMAGES_BUCKET"`
	IDBackImagesBucket        string `mapstructure:"ID_BACK_IMAGES_BUCKET"`
	WalletSecret              string `mapstructure:"LITPAD_WALLET_SECRET"`
	PGAdminPassword           string `mapstructure:"PGADMIN_PASSWORD"`
	ICPWalletIp               string `mapstructure:"LITPAD_WALLET_IP"`
	ICPPrivateKey             string `mapstructure:"ICP_PRIVATE_KEY"`
	ICPPublicKey              string `mapstructure:"ICP_PUBLIC_KEY"`
	CloudinaryCloudName       string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryApiKey          string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecret       string `mapstructure:"CLOUDINARY_API_SECRET"`
	RedisUrl                  string `mapstructure:"REDIS_URL"`
	ReminderCronHours         uint   `mapstructure:"REMINDER_CRON_HOURS"`
	AppScheme                 string `mapstructure:"APP_SCHEME"`
}

func GetConfig() (config Config) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "." // Default to current directory if not set
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	var err error
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
	config.SecretKeyByte = []byte(config.SecretKey)
	return
}
