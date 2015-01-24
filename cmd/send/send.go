package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
)

type Message struct {
	To      []*mail.Address
	From    *mail.Address
	Subject string
	Headers []string
	Content string
}

func NewMessage(to, from, subject string, headers []string, content string) *Message {
	m := new(Message)
	m.To, _ = mail.ParseAddressList(to)
	m.From, _ = mail.ParseAddress(from)
	m.Subject = subject
	m.Headers = headers
	m.Content = content

	return m
}

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

func (m *Message) allHeaders() string {
	var output string

	for _, h := range m.Headers {
		output += h + "\n"
	}

	return output
}

type SMTPServer struct {
	name     string
	username string
	password string
	address  string
	port     int
	tlsB     bool
}

func staticSMTPServer() *SMTPServer {
	s := new(SMTPServer)
	s.name = "smtp.gmail.com"
	s.username = "mr.k.frenata@gmail.com"
	s.password = "gophergala"
	s.address = "smtp.gmail.com"
	s.port = 587
	s.tlsB = true

	return s
}

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

func sendSMTP(server *SMTPServer, msg *Message) error {
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
	s := staticSMTPServer()
	m := NewMessage("danweasel@gmail.com, dan.weasel@gmail.com", "mr.k.frenata@gmail.com", "Is Anybody Out There?", []string{"Reply-To: Frenata <mr.k.frenata@gmail.com>", "Content-Type: text/plain"}, "Ground control to Major Tom!")
	fmt.Println(sendSMTP(s, m))
}
