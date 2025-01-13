package xmail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// Config represents the email server configuration.
type Config struct {
	SMTPHost    string
	SMTPPort    int
	Username    string
	Password    string
	FromAddress string
}

// Email represents the details of an email to be sent.
type Email struct {
	To           []string
	Subject      string
	TemplateFile string
	Data         any
}

// XMail is the main structure for the xmail package.
type XMail struct {
	config Config
}

// New creates a new instance of XMail.
func New(config Config) XMail {
	return XMail{config: config}
}

// Send sends a single email.
func (x XMail) Send(email Email) error {
	// Parse the template file.
	tmpl, err := template.ParseFiles(email.TemplateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	// Render the template with the provided data.
	var body bytes.Buffer
	headers := fmt.Sprintf("Subject: %s\r\n", email.Subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	body.WriteString(headers)

	if err := tmpl.Execute(&body, email.Data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Set up authentication.
	auth := smtp.PlainAuth("", x.config.Username, x.config.Password, x.config.SMTPHost)

	// Send the email.
	addr := fmt.Sprintf("%s:%d", x.config.SMTPHost, x.config.SMTPPort)
	if err := smtp.SendMail(addr, auth, x.config.FromAddress, email.To, body.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Example usage:
// func main() {
// 	config := Config{
// 		SMTPHost:    "smtp.example.com",
// 		SMTPPort:    587,
// 		Username:    "your_username",
// 		Password:    "your_password",
// 		FromAddress: "you@example.com",
// 	}
//
// 	xmail := New(config)
//
// 	email := Email{
// 		To:          []string{"recipient@example.com"},
// 		Subject:     "Test Email",
// 		TemplateFile: "example.tpl",
// 		Data:        map[string]string{"Name": "John Doe"},
// 	}
//
// 	if err := xmail.Send(email); err != nil {
// 		fmt.Println("Failed to send email:", err)
// 	} else {
// 		fmt.Println("Email sent successfully!")
// 	}
// }
