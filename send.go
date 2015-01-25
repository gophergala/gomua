package gomua

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const configLocation = "/.gomua/send.cfg"

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
		return nil, errors.New("Missing ~/.gomua/send.cfg file.")
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

	if s.name == "" || s.username == "" || s.password == "" || s.address == "" || s.port == 0 {
		return nil, errors.New("Incorrect ~/.mua/send.cfg file.")
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
func sendSMTP(server *SMTPServer, msg *Message) error {
	// connect to SMTP server
	var c *smtp.Client
	c, err := connectSMTP(server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Set sender and receiver
	var from string
	if msg.from == nil {
		from = server.username
	} else {
		from = msg.from.Address
	}
	if err := c.Mail(from); err != nil {
		log.Fatal(err)
	}

	for _, t := range msg.to {
		if err := c.Rcpt(t.Address); err != nil {
			log.Fatal(err)
		}
	}

	// Send email body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	msg.SetDate(time.Now().Local())
	_, err = fmt.Fprint(wc, msg.Headers())
	_, err = fmt.Fprint(wc, msg.String())
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Send opens a new SMTP server connection from the config file and sends a message.
func Send(msg *Message) {
	// Look for a SMTPServer configuration file in ~/.gomua/send.cfg
	u, _ := user.Current()
	srv, err := NewSMTPServer(u.HomeDir + configLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nSending...")

	if sendSMTP(srv, msg) == nil {
		fmt.Println("Message Sent")
	}
}
