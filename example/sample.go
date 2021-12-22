package example

import (
	"fmt"

	"github.com/dsa0x/sicher/sicher"
)

type Config struct {
	Port                    string `required:"true" envconfig:"PORT"`
	MongoDbURI              string `required:"true" envconfig:"MONGO_DB_URI"`
	MongoDbName             string `required:"true" envconfig:"MONGO_DB_NAME"`
	JWTSecret               string `required:"true" envconfig:"JWT_SECRET"`
	S3Bucket                string `required:"true" envconfig:"S3_BUCKET"`
	AWSAccessKey            string `required:"true" envconfig:"AWS_ACCESS_KEY"`
	AWSSecretKey            string `required:"true" envconfig:"AWS_SECRET_KEY"`
	AWSRegion               string `required:"true" envconfig:"AWS_REGION"`
	AWSLogGroupName         string `required:"true" envconfig:"AWS_CLOUDWATCHLOGS_GROUP_NAME"`
	AWSLogStreamName        string `required:"true" envconfig:"AWS_CLOUDWATCHLOGS_STREAM_NAME"`
	Environment             string `required:"true" envconfig:"ENVIRONMENT"`
	SecretReferralID        string `required:"true" envconfig:"SECRET_REFERRAL_ID"`
	SendGridAPIKey          string `required:"false" envconfig:"SENDGRID_API_KEY"`
	SenderEmail             string `required:"true" envconfig:"SENDER_EMAIL"`
	SenderName              string `required:"true" envconfig:"SENDER_NAME"`
	FlutterwaveSecretKey    string `required:"true" envconfig:"FLUTTERWAVE_SECRET_KEY"`
	FlutterwaveAcctEndpoint string `required:"true" envconfig:"FLUTTERWAVE_ACCT_ENDPOINT"`
}

func Configure() {
	var config Config

	env := "development"
	path := "."
	s := sicher.New(env, path)
	err := s.LoadEnv("", &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.SenderName)
}
