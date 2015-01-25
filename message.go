package gomua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"time"
)

// Message contains critical data necessary for readding/sending a message.
type Message struct {
	to      []*mail.Address
	from    *mail.Address
	date    time.Time
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

// Strings pretty prints Message, with just the standard headers.
func (m *Message) String() string {
	return fmt.Sprintf(
		"From: %s\nTo: %s\nDate: %v\nSubject: %s\n\n%s\n",
		m.From(), m.To(), m.Date(), m.Subject(), m.Content())
}

func (m *Message) Subject() string        { return m.subject }
func (m *Message) Content() string        { return m.content }
func (m *Message) Date() string           { return m.date.Format(time.RFC822) }
func (m *Message) SetDate(date time.Time) { m.date = date }

func (m *Message) From() string {
	if m.from != nil {
		return m.from.String()
	} else {
		return ""
	}
}

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

// Save writes the message to a file.
func (m *Message) Save(file string) error {
	b := bytes.NewBufferString(m.String()).Bytes()
	return ioutil.WriteFile(file, b, 0660)
}
