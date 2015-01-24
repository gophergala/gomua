package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"time"
)

type Message struct {
	To         []*mail.Address
	From       *mail.Address
	Subject    string
	Date       time.Time
	ReplyID    string
	References []string
	Content    string
}

func NewMessage(to string, from string, subject string, date time.Time, replyID string, references []string, content string) (m Message) {
	m.To, _ = mail.ParseAddressList(to)
	m.From, _ = mail.ParseAddress(from)
	m.Subject = subject
	m.Date = date
	m.ReplyID = replyID
	m.References = references
	m.Content = content

	return m
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

func sendSMTP(server *SMTPServer, msg Message) error {
	// connect to SMTP server
	var c *smtp.Client
	c, err := connectSMTP(server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Set sender and receiver
	if err := c.Mail("mr.k.frenata@gmail.com"); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt("danweasel@gmail.com"); err != nil {
		log.Fatal(err)
	}

	// Send email body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprint(wc, "To: Andrew <danweasel@gmail.com>\nFrom: Frenata <mr.k.frenata@gmail.com>\nSubject: testing message from Go\n\nGround control to Major Tom!")
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
	m := NewMessage("danweasel@gmail.com", "mr.k.frenata@gmail.com", "Test Message", time.Now(), "", nil, "This is test message content.")
	fmt.Println(sendSMTP(s, m))
}
