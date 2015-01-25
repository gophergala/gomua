package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const configLocation = "/.gomua/send.cfg"

// Message contains critical data necessary for sending a message.
type Message struct {
	to      []*mail.Address
	from    *mail.Address
	subject string
	headers []string
	content string
}

// NewMessage takes strings that specify the recipients, the sender, the subject, a slice of other email header values, and the content of an email.
func NewMessage(to, from, subject string, headers []string, content string) *Message {
	m := new(Message)
	m.to, _ = mail.ParseAddressList(to)
	m.from, _ = mail.ParseAddress(from)
	m.subject = subject
	m.headers = headers
	m.content = content

	return m
}

func (m *Message) From() string    { return m.from.String() }
func (m *Message) Subject() string { return m.subject }
func (m *Message) Content() string { return m.content }

// Outputs the list of parsed addresses back into a single string for appending to the email.
func (m *Message) To() (output string) {
	for i, f := range m.to {
		output += f.String()
		if len(m.to)-1 > i {
			output += ","
		}
	}
	return output
}

// Writes the slice of headers one on each line.
func (m *Message) Headers() (output string) {
	for _, h := range m.headers {
		output += h + "\n"
	}
	return output
}

func (m *Message) AddHeader(header string) {
	m.headers = append(m.headers, header)
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
	if err := c.Mail(msg.from.Address); err != nil {
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
	msg.AddHeader(fmt.Sprintf("Date: %s%s", time.Now().Local().Format(time.RFC822), "\r\n"))
	_, err = fmt.Fprintf(wc, "To: %v\nFrom: %v\nSubject: %v\n%v\n\n%v", msg.To(), msg.From(), msg.Subject(), msg.Headers(), msg.Content())
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// WriteMessage interactively prompts the user for an email to send.
func WriteMessage(r io.Reader) *Message {
	cli := bufio.NewScanner(r)

	fmt.Print("To: ")
	cli.Scan()
	to := cli.Text()
	fmt.Print("From: ")
	cli.Scan()
	from := cli.Text()
	fmt.Print("Subject: ")
	cli.Scan()
	subject := cli.Text()
	fmt.Print("Content: (Enter SEND to finish adding content and send the email.\n")

	var content string
	for {
		cli.Scan()
		line := cli.Text()
		if line == "SEND" {
			break
		} else {
			content += line + "\n"
		}
	}

	msg := NewMessage(to, from, subject, nil, content)
	msg.AddHeader("Content-Type: text/plain; charset=UTF-8")
	return msg
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

func main() {
	// write a message on the stdin and send it
	Send(WriteMessage(os.Stdin))
}
