package shemail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

const (
	fromAddress = "test@example.co.kr"
	smtpServer  = "smtp.mailplug.co.kr:465"
)

var (
	smtpPassword = ""
)

func SetSmtpPassword() {
	smtpPassword = os.Getenv("SMTP_SERVER_PASSWORD")
}

// Send email to address
func SendEmail(address string, subject string, body string) error {
	from := mail.Address{Name: "", Address: fromAddress}
	to := mail.Address{Name: "", Address: address}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host, port, err := net.SplitHostPort(smtpServer)
	if err != nil {
		return err
	}
	fmt.Printf("host: %s, port: %s\n", host, port)

	auth := smtp.PlainAuth("", fromAddress, smtpPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}
	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
