package email

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/COTBU/notifier/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type EmailClient interface {
	Credentials() *Credentials
	SetDestination([]string) EmailClient
	SetSubject(string) EmailClient
	SendRich(string, any) error
}

// Credentials struct.
type Credentials struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// Email struct
// https://medium.com/@dhanushgopinath/sending-html-emails-using-templates-in-golang-9e953ca32f3d
type client struct {
	config      *config.Config
	credentials *Credentials
	smtpClient  *smtp.Client
}

//go:embed templates/_header.html
var header string

//go:embed templates/_footer.html
var footer string

const moscowSecondsUTCOffset = 3 * 60 * 60

var moscowLocation = func() *time.Location {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		location = time.FixedZone("UTC+3", moscowSecondsUTCOffset)
	}
	return location
}()

// NewClient func
func NewClient(config *config.Config) (EmailClient, error) {
	s := &client{
		config:      config,
		credentials: &Credentials{From: config.Email.User},
	}

	if err := s.newClient(); err != nil {
		return nil, err
	}

	return s, nil
}

func (c *client) newClient() (err error) {
	servername := fmt.Sprintf("%s:%d", c.config.Email.Host, c.config.Email.Port)

	c.smtpClient, err = smtp.Dial(servername)
	if err != nil {
		return err
	}
	if runtime.GOOS == "linux" {
		// TLS config
		host, _, _ := net.SplitHostPort(servername)
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			ServerName:         host,
		}
		if err = c.smtpClient.StartTLS(tlsconfig); err != nil {
			return err
		}
		// Auth
		auth := smtp.PlainAuth(
			"",
			c.config.Email.User,
			c.config.Email.Password,
			c.config.Email.Host,
		)
		return c.smtpClient.Auth(auth)
	}

	return nil
}

func (c *client) Credentials() *Credentials {
	return c.credentials
}

func (c *client) SetDestination(to []string) EmailClient {
	c.credentials.To = to
	return c
}

func (c *client) SetSubject(subject string) EmailClient {
	c.credentials.Subject = subject
	return c
}

// TemplateData struct.
type TemplateData struct {
	Data any
	IP   string
}

// SendRich - отправляет письмо с разметкой.
// Принимает на вход строку шаблона в виде MD и данные для шаблона.
// Данные шаблона обогащаются данными об IP адресе сервера.
// Возвращает ошибку в случае если невозможно распарсить шаблон или отправить письмо.
func (c *client) SendRich(templateBody string, data any) (err error) {
	emailData := TemplateData{data, c.config.Service.Address}
	c.credentials.Body, err = c.execTemplate(emailData, templateBody)
	if err != nil {
		err = fmt.Errorf(
			"error parsing template: %w,\nemailData: %+v,\ntemplate: %s",
			err,
			emailData,
			templateBody,
		)
		sentry.CaptureException(err)
		return err
	}

	return c.send()
}

func (c *client) send() error {
	header := map[string]string{
		"To":                        strings.Join(c.credentials.To, ","),
		"From":                      c.credentials.From,
		"Subject":                   c.credentials.Subject,
		"MIME-Version":              "1.0",
		"Content-Type":              "text/html; charset=\"utf-8\"",
		"Content-Transfer-Encoding": "base64",
	}
	sb := strings.Builder{}
	for title, value := range header {
		_, _ = sb.WriteString(fmt.Sprintf("%s: %s\r\n", title, value))
	}
	_, _ = sb.WriteString("\r\n" + base64.StdEncoding.EncodeToString([]byte(c.credentials.Body)))

	// To && From
	if err := c.smtpClient.Mail(c.credentials.From); err != nil {
		return err
	}
	for _, t := range c.credentials.To {
		if err := c.smtpClient.Rcpt(t); err != nil {
			return err
		}
	}
	// Data
	w, err := c.smtpClient.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(sb.String())); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return c.smtpClient.Quit()
}

func (c *client) execTemplate(data any, templ string) (res string, err error) {
	funcs := template.FuncMap{
		"getRequestFiles": GetRequestFiles,
		"formatDate": func(date time.Time) string {
			return date.In(moscowLocation).Format("02.01.2006 15:04")
		},
	}
	var buf bytes.Buffer
	t, err := template.New("title").Funcs(funcs).Parse(templ)
	if err != nil {
		return
	}
	if err = t.Execute(&buf, data); err != nil {
		return
	}

	buffer := new(bytes.Buffer)
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	if err := md.Convert(buf.Bytes(), buffer); err != nil {
		return "", err
	}

	res = header + buffer.String() + footer

	return res, nil
}

func GetRequestFiles(sType uint) string {
	var files []string
	if (sType & 1) != 0 {
		files = append(files, "TXT")
	}
	if (sType & 2) != 0 {
		files = append(files, "PDF")
	}
	if (sType & 4) != 0 {
		files = append(files, "Excel")
	}
	return strings.Join(files, " + ")
}
