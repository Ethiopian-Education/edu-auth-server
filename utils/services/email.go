package services

import (
	"bytes"
	"html/template"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/go-mail/mail"
)

// USING GO-MAIL package that's built OVER THE native smtp/template go package

func SendEmailMessage(to []string, from string, subject string, body string) error {
	m  := mail.NewMessage()

	m.SetHeader("From",from)

	m.SetHeader("To", to...)

	m.SetAddressHeader("Cc", "abemelekmila@gmail.com", "Abemelek")

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	dial := mail.NewDialer(config.SMTP_HOST, 465, config.SMTP_USERNAME, config.SMTP_PASSWORD)

	// Now send the message
	if err := dial.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func ParseHtmlTemplate(template_filename string, data interface{}) (string, error) {
	
	t, err := template.ParseFiles(template_filename)
	if err != nil {
		return "" , err
	}

	buf := new(bytes.Buffer) // memory allocation using new() operator...

	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	// Parsing finished ...

	return buf.String(), nil

}