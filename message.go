package gomua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
)

type Message struct {
	mail.Message
	Content  string
	Filename string
}

func (m *Message) store() {
	b, err := ioutil.ReadAll(m.Message.Body)
	if err != nil {
		log.Fatal(err)
	}
	m.Content = string(b)
}

// ReadMessage embeds a mail.Message inside the gomua.Message and stores the Body content
func ReadMessage(msg *mail.Message) *Message {
	m := new(Message)
	m.Message = *msg
	m.store()

	return m
}

func (m *Message) Seen() {
	s := strings.Split(m.Filename, ":2")
	if len(s) != 2 {
		log.Fatal(fmt.Errorf("filename %s does not contain ':2'", m.Filename))
	}
	name := s[0]
	flags := s[1]

	if !strings.Contains(flags, "S") {
		if flags[len(flags)-1] != ',' {
			flags += ","
		}
		flags += "S"

		newname := name + ":2" + flags
		os.Rename(m.Filename, newname)
	}
}

// WriteMessage interactively prompts the user for an email to send.
func WriteMessage(r io.Reader) *mail.Message {
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
	content := WriteContent(r)

	msg := "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += fmt.Sprintf(
		"To: %v\r\nFrom: %v\r\nSubject: %v\r\n\r\n%v",
		to, from, subject, content)

	m, err := mail.ReadMessage(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func WriteContent(r io.Reader) string {
	cli := bufio.NewScanner(r)
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
	return content
}

func Save(file string, m string) error {
	b := bytes.NewBufferString(m).Bytes()
	ioutil.WriteFile(file, b, 0600)
	return nil
}
