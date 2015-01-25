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

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func scanMailDir(dir string) (msgs []gomua.Mail) {
	mails := gomua.Scan(dir)

	// Embed mail.Message inside gomua.Message
	for _, m := range mails {
		msg := gomua.ReadMessage(m)
		msgs = append(msgs, msg)
	}

	return msgs
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMailList(msgs []gomua.Mail, w io.Writer) {
	for i, m := range msgs {
		fmt.Fprintf(w, "%d. %s\n", i+1, m.Summary())
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
func replyMessage(old *gomua.Message) (reply *mail.Message) {
	oldid := old.Header.Get("Message-ID")
	oldref := old.Header.Get("References")

	to := "To: " + old.Header.Get("From") + "\r\n"
	// TODO: pull from config?
	from := "From: mr.k.frenata@gmail.com\r\n"

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

// user input loop
func input(mails []gomua.Mail, exit chan bool) {
	cli := bufio.NewScanner(os.Stdin)
	for {
		cli.Scan()
		input := cli.Text()

		switch {
		case input == "main":
			viewMailList(mails, os.Stdout)
		case strings.HasPrefix(input, "reply"):
			num, err := strconv.Atoi(strings.TrimPrefix(input, "reply "))
			if err != nil {
				fmt.Println(err)
			}
			reply := replyMessage(mails[num-1].(*gomua.Message))
			gomua.Send(reply)
			//viewMessage(reply)
			//gomua.Send(&reply.Message)
		case input == "exit", input == "x", input == "quit", input == "q":
			exit <- true
		case strings.ContainsAny(input, "01234566789"):
			num, _ := strconv.Atoi(input)
			if num <= len(mails) && num > 0 {
				viewMail(mails[num-1], os.Stdout)
			}
		}

	}

	exit <- true
}

func main() {
	// TODO: read from config? or command line arg?
	const dir string = "./testmaildir"

	msgs := scanMailDir(dir)
	viewMailList(msgs, os.Stdout)

	exit := make(chan bool, 1)
	go input(msgs, exit)

	<-exit
	os.Exit(2)
}
