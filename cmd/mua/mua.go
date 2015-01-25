package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/gophergala/gomua"
)

// client handles common data as a user navigates the MUA.
// messages is the list of all mail (in the current folder?)
// current is the currently selected Mail
// displayN is the # of Mail to display on the screen at one time.
// user is the user's email address, for sending.
type client struct {
	messages []gomua.Mail
	current  gomua.Mail
	displayN int
	user     string
}

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func (c *client) scanMailDir(dir string) {
	var msgs []gomua.Mail
	mails := gomua.Scan(dir)

	// Embed mail.Message inside gomua.Message
	for _, m := range mails {
		msg := gomua.ReadMessage(m)
		msgs = append(msgs, msg)
	}

	c.messages = msgs
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMailList(msgs []gomua.Mail, start int, w io.Writer) {
	for i, m := range msgs {
		fmt.Fprintf(w, "%d. %s\n", i+start+1, m.Summary())
	}
}

// prints a single mail message to the screen
func viewMail(msg gomua.Mail, w io.Writer) {
	fmt.Fprint(w, msg)
}

// (to be invoked from viewMessage with a keypress, most likely)
// create a new message:
// put old's Message-ID into reply's In-Reply-To and References headers
// put old's content as quoted material in reply
// put old's From as the reply's To
// put old's Subject in reply's subject as "RE: Subject" (but allow user to edit? -- later)
// prompt user for content
// send reply
func replyMessage(old *gomua.Message, user string) (reply *mail.Message) {
	oldid := old.Header.Get("Message-ID")
	oldref := old.Header.Get("References")

	to := "To: " + old.Header.Get("From") + "\r\n"
	from := fmt.Sprintf("From: %s\r\n", user)

	subject := fmt.Sprintf("Subject: RE: %s\r\n", old.Header.Get("Subject"))
	inreplyto := fmt.Sprintf("In-Reply-To: %s\r\n", oldid)
	references := fmt.Sprintf("References: %s%s\r\n", oldref, oldid)

	content := "\r\n\r\n" + gomua.WriteContent(os.Stdin)
	quote := bufio.NewScanner(strings.NewReader(old.Content))
	for quote.Scan() {
		content += "\n" + "> " + quote.Text()
	}

	buf := bytes.NewBufferString(inreplyto)
	buf.WriteString(references)
	buf.WriteString(to)
	buf.WriteString(from)
	buf.WriteString(subject)
	buf.WriteString(content)

	reply, err := mail.ReadMessage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	return reply
}

// adds ANSI colors to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}

// helper func to check that view doesn't overflow []. Returns end.
func (c *client) printList(start, end int) (newstart int, newend int) {
	if end = start + c.displayN; end > len(c.messages) {
		end = len(c.messages)
	}
	m := c.messages[start:end]
	viewMailList(m, start, os.Stdout)
	start = end
	return start, end
}

// user input loop
func (c *client) input(exit chan bool) {
	var start, end int
	cli := bufio.NewScanner(os.Stdin)
	for {
		cli.Scan()
		input := cli.Text()

		switch {
		case input == "main", input == "view":
			start = 0
			start, end = c.printList(start, end)
		case input == "more":
			start, end = c.printList(start, end)
		case strings.HasPrefix(input, "reply"):
			num, err := strconv.Atoi(strings.TrimPrefix(input, "reply "))
			if err != nil {
				fmt.Println(err)
			}
			reply := replyMessage(c.messages[num-1].(*gomua.Message), c.user)
			gomua.Send(reply)
		case input == "exit", input == "x", input == "quit", input == "q":
			exit <- true
		case strings.ContainsAny(input, "01234566789"):
			num, _ := strconv.Atoi(input)
			if num <= len(c.messages) && num > 0 {
				viewMail(c.messages[num-1], os.Stdout)
			}
		}

	}

	exit <- true
}

func main() {
	// TODO: read from config? or command line arg?
	const dir string = "./testmaildir"

	client := &client{displayN: 25, user: "Frenata <mr.k.frenata@gmail.com>"}
	client.scanMailDir(dir)

	exit := make(chan bool, 1)
	go client.input(exit)

	<-exit
	os.Exit(2)
}
