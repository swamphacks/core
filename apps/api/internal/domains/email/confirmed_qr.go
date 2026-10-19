package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

func (s *EmailService) SendConfirmedQREmail(
	ctx context.Context, userID uuid.UUID, recipient, name string,
) (string, error) {
	address, err := mail.ParseAddress(recipient)
	if err != nil || address.Address != recipient {
		return "", fmt.Errorf("invalid recipient address")
	}
	t, err := template.ParseFiles(filepath.Join(
		s.config.EmailTemplateDirectory, "ConfirmedCheckInEmail.html",
	))
	if err != nil {
		return "", err
	}
	var htmlBody bytes.Buffer
	if err := t.Execute(&htmlBody, struct{ Name string }{name}); err != nil {
		return "", err
	}
	png, err := qrcode.Encode("IDENT::"+userID.String(), qrcode.Highest, 512)
	if err != nil {
		return "", err
	}

	var message bytes.Buffer
	mixed := multipart.NewWriter(&message)
	subject := "IMPORTANT: Your SwampHacks XII QR Code, Check-In Details, and Reminders"
	fmt.Fprintf(&message, "From: SwampHacks <contact@swamphacks.com>\r\n")
	fmt.Fprintf(&message, "To: %s\r\n", address.String())
	fmt.Fprintf(&message, "Reply-To: contact@swamphacks.com\r\n")
	fmt.Fprintf(&message, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", subject))
	fmt.Fprintf(&message, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&message, "Content-Type: multipart/mixed; boundary=%q\r\n\r\n", mixed.Boundary())

	addPart := func(contentType, disposition string, data []byte) error {
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", contentType)
		header.Set("Content-Transfer-Encoding", "base64")
		if disposition != "" {
			header.Set("Content-Disposition", disposition)
		}
		part, err := mixed.CreatePart(header)
		if err != nil {
			return err
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		for len(encoded) > 0 {
			n := 76
			if len(encoded) < n {
				n = len(encoded)
			}
			if _, err := fmt.Fprintf(part, "%s\r\n", encoded[:n]); err != nil {
				return err
			}
			encoded = encoded[n:]
		}
		return nil
	}
	if err := addPart("text/html; charset=utf-8", "", htmlBody.Bytes()); err != nil {
		return "", err
	}
	if err := addPart("image/png",
		`attachment; filename="SwampHacks-XII-Check-In-QR.png"`, png); err != nil {
		return "", err
	}
	if err := mixed.Close(); err != nil {
		return "", err
	}
	return s.SESClient.SendRawMessage(ctx, recipient, message.Bytes())
}
