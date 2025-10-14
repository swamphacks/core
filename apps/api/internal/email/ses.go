package email

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/rs/zerolog"
)

type SESClient struct {
	client *ses.Client
	logger zerolog.Logger
}

func NewSESClient(accountId, accessKey, secretkey, region string, logger zerolog.Logger) *SESClient {
	// Initialize the SES client here with appropriate configuration
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretkey, "")),
		awsConfig.WithRegion(region),
	)
	if err != nil {
		logger.Err(err).Msg("Failed to load AWS configuration")
		return nil
	}

	client := ses.NewFromConfig(cfg)

	return &SESClient{
		client: client,
		logger: logger.With().Str("component", "ses_client").Logger(),
	}
}

func (c *SESClient) SendEmail(to []string, from, subject string, body string) error {
	_, err := c.client.SendEmail(context.TODO(), &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: to,
		},
		Source: &from,
		Message: &types.Message{
			Subject: &types.Content{
				Data: &subject,
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: &body,
				},
			},
		},
	})
	if err != nil {
		c.logger.Err(err).Msg("Failed to send email")
		return err
	}

	return nil
}

func (c *SESClient) SendHTMLEmail(to []string, from, subject string, htmlBody string) error {
	_, err := c.client.SendEmail(context.TODO(), &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: to,
		},
		Source: &from,
		Message: &types.Message{
			Subject: &types.Content{
				Data: &subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: &htmlBody,
				},
			},
		},
	})
	if err != nil {
		c.logger.Err(err).Msg("Failed to send HTML email")
		return err
	}

	return nil
}
