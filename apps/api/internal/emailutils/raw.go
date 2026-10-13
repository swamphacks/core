package emailutils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func (c *SESClient) SendRawMessage(ctx context.Context, recipient string, message []byte) (string, error) {
	source := "SwampHacks <contact@swamphacks.com>"
	result, err := c.client.SendRawEmail(ctx, &ses.SendRawEmailInput{
		Source:       &source,
		Destinations: []string{recipient},
		RawMessage:   &types.RawMessage{Data: message},
	})
	if err != nil {
		return "", err
	}
	if result.MessageId == nil {
		return "", nil
	}
	return *result.MessageId, nil
}
