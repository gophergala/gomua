package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os/user"
	"strconv"
	"strings"
)

// Message contains critical data necessary for sending a message.
type Message struct {
	To      []*mail.Address
	From    *mail.Address
	Subject string
	Headers []string
	Content string
}

// NewMessage takes strings that specify the recipients, the sender, the subject, a slice of other email header values, and the content of an email.
func NewMessage(to, from, subject string, headers []string, content string) *Message {
	m := new(Message)
	m.To, _ = mail.ParseAddressList(to)
	m.From, _ = mail.ParseAddress(from)
	m.Subject = subject
	m.Headers = headers
	m.Content = content

	return m
}

// Outputs the list of parsed addresses back into a single string for appending to the email.
func (m *Message) allTo() string {
	var output string

	for i, f := range m.To {
		output += f.String()
		if len(m.To)-1 > i {
			output += ","
		}
	}

	return output
}

// Writes the slice of headers one on each line.
func (m *Message) allHeaders() string {
	var output string

	for _, h := range m.Headers {
		output += h + "\n"
	}

	return output
}

// SMTPServer describes a connection to an SMTP server for sending mail.
type SMTPServer struct {
	name     string
	username string
	password string
	address  string
	port     int
	tlsB     bool
}

// NewSMTPServer reads from a configuration file, and returns a new SMTPServer struct ready to use.
func NewSMTPServer(filename string) (*SMTPServer, error) {
	s := new(SMTPServer)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	for _, l := range lines {
		switch {
		case strings.HasPrefix(l, "Name="):
			s.name = strings.TrimPrefix(l, "Name=")
		case strings.HasPrefix(l, "Username="):
			s.username = strings.TrimPrefix(l, "Username=")
		case strings.HasPrefix(l, "Password="):
			s.password = strings.TrimPrefix(l, "Password=")
		case strings.HasPrefix(l, "Address="):
			s.address = strings.TrimPrefix(l, "Address=")
		case strings.HasPrefix(l, "Port="):
			s.port, _ = strconv.Atoi(strings.TrimPrefix(l, "Port="))
		case strings.HasPrefix(l, "TLS="):
			str := strings.TrimPrefix(l, "TLS=")
			if str == "true" {
				s.tlsB = true
			} else {
				s.tlsB = false
			}

		}
	}
	return s, nil
}

// Connects and authenticates to an SMTPServer, returns a client connection ready to write.
// This client *must be Quit()ed after finished using, preferably with defer.
func connectSMTP(s *SMTPServer) (*smtp.Client, error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return nil, err
	}

	if err := c.StartTLS(&tls.Config{ServerName: s.name}); err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.address)
	if err := c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}

// SendSMTP takes a SMTP server and a message, connects to the server, sends the message, and quits the connection to the server.
func SendSMTP(server *SMTPServer, msg *Message) error {
	// connect to SMTP server
	var c *smtp.Client
	c, err := connectSMTP(server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Set sender and receiver
	if err := c.Mail(msg.From.Address); err != nil {
		log.Fatal(err)
	}

	for _, t := range msg.To {
		if err := c.Rcpt(t.Address); err != nil {
			log.Fatal(err)
		}
	}

	// Send email body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "To: %v\nFrom: %v\nSubject: %v\n%v\n\n%v", msg.allTo(), msg.From, msg.Subject, msg.allHeaders(), msg.Content)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	// Look for a SMTPServer configuration file in ~/.mua/send.cfg
	u, _ := user.Current()
	s, err := NewSMTPServer(u.HomeDir + "/.mua/send.cfg")
	if err != nil {
		fmt.Println(err)
		return
	}

	m := NewMessage("danweasel@gmail.com, dan.weasel@gmail.com", "mr.k.frenata@gmail.com", "Is Anybody Out There?", []string{"Reply-To: Frenata <mr.k.frenata@gmail.com>", "Content-Type: text/plain"}, "Ground control to Major Tom!")
	if SendSMTP(s, m) == nil {
		fmt.Println("Message sent sucessfully!")
	}
}
