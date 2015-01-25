package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/gophergala/gomua"
)

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func scanMailDir(dir string) (mail []*mail.Message) {

	return mail
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMail(mail []*mail.Message, w io.Writer) {
	for i, m := range mail {
		subject := m.Header.Get("Subject")
		from := m.Header.Get("From")
		fmt.Fprintf(w, "%d. %s from %s\n", i+1, color(subject, "31"), color(from, "33"))
	}
}

// prints a single mail message to the screen
func viewMessage(msg *mail.Message) {
	fmt.Printf("From: %v\n", msg.Header.Get("From"))
	fmt.Printf("To: %v\n", msg.Header.Get("To"))
	fmt.Printf("Date: %v\n", msg.Header.Get("Date"))
	fmt.Printf("Subject: %v\n", msg.Header.Get("Subject"))
	body, _ := ioutil.ReadAll(msg.Body)
	fmt.Printf("\n%v", string(body))
}

// (to be invoked from viewMessage with a keypress, most likely)
// create a new message:
// put old's Message-ID into reply's In-Reply-To and References headers
// put old's content as quoted material in reply
// put old's From as the reply's To
// put old's Subject in reply's subject as "RE: Subject" (but allow user to edit? -- later)
// prompt user for content
// send reply
func replyMessage(old *mail.Message) (reply *mail.Message) {
	//id := "121212"                             // should be parsed from old by hp.go eventually
	//oldref := ""                               // ditto
	//from := "Frenata <mr.k.frenata@gmail.com>" // mua.cfg for sending email address?

	//subject := fmt.Sprintf("RE: %s", old.Header.Get("Subject"))
	//inreplyto := fmt.Sprintf("In-Reply-To: %s", id)
	//references := fmt.Sprintf("References: %s%s", oldref, id)
	//headers := []string{inreplyto, references}

	content := gomua.WriteContent(os.Stdin)
	quote := bufio.NewScanner(old.Body)
	for quote.Scan() {
		content += "\n" + "> " + quote.Text()
	}

	//reply = gomua.NewMessage(old.From(), from, subject, headers, content)
	return reply
}

// adds ANSI colors to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}

// user input loop
func input(mails []*mail.Message, exit chan bool) {
	cli := bufio.NewScanner(os.Stdin)
	for {
		cli.Scan()
		input := cli.Text()

		switch {
		case input == "main":
			viewMail(mails, os.Stdout)
		case strings.HasPrefix(input, "reply"):
			num, err := strconv.Atoi(strings.TrimPrefix(input, "reply "))
			if err != nil {
				fmt.Println(err)
			}
			reply := replyMessage(mails[num-1])
			gomua.Send(reply)
		case input == "exit", input == "x", input == "quit", input == "q":
			exit <- true
		case strings.ContainsAny(input, "01234566789"):
			num, _ := strconv.Atoi(input)
			if num <= len(mails) && num > 0 {
				viewMessage(mails[num-1])
			}
		}

	}

	exit <- true
}

func main() {
	mails := fakeMailDir()
	viewMail(mails, os.Stdout)

	exit := make(chan bool, 1)
	go input(mails, exit)

	<-exit
	os.Exit(2)
}

// a list of static Messages for testing view
func fakeMailDir() (mails []*mail.Message) {
	b, err := ioutil.ReadFile("m1.mail")
	if err != nil {
		log.Fatal(err)
	}

	m1, err := mail.ReadMessage(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	mails = append(mails, m1)
	return mails
}
