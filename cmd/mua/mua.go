package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala/gomua"
)

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func scanMailDir(dir string) (mail []*gomua.Message) {

	return mail
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMail(mail []*gomua.Message, w io.Writer) {
	for i, m := range mail {
		fmt.Fprintf(w, "%d. %s from %s\n", i+1, color(m.Subject(), "31"), color(m.From(), "33"))
	}
}

// prints a single mail message to the screen
func viewMessage(msg *gomua.Message) {
	fmt.Println(msg)
}

// (to be invoked from viewMessage with a keypress, most likely)
// create a new message:
// put old's Message-ID into reply's In-Reply-To and References headers
// put old's content as quoted material in reply
// put old's From as the reply's To
// put old's Subject in reply's subject as "RE: Subject" (but allow user to edit? -- later)
// prompt user for content
// send reply
func replyMessage(old *gomua.Message) (reply *gomua.Message) {
	id := "121212"                             // should be parsed from old by hp.go eventually
	oldref := ""                               // ditto
	from := "Frenata <mr.k.frenata@gmail.com>" // mua.cfg for sending email address?

	subject := fmt.Sprintf("RE: %s", old.Subject())
	inreplyto := fmt.Sprintf("In-Reply-To: %s", id)
	references := fmt.Sprintf("References: %s%s", oldref, id)
	headers := []string{inreplyto, references}

	content := gomua.WriteContent(os.Stdin)

	quote := bufio.NewScanner(strings.NewReader(old.Content()))
	for quote.Scan() {
		content += "\n" + "> " + quote.Text()
	}

	reply = gomua.NewMessage(old.From(), from, subject, headers, content)
	return reply
}

// adds ANSI colors to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}

// user input loop
func input(mail []*gomua.Message, exit chan bool) {
	cli := bufio.NewScanner(os.Stdin)
	for {
		cli.Scan()
		input := cli.Text()

		switch {
		case input == "main":
			viewMail(mail, os.Stdout)
		case strings.HasPrefix(input, "reply"):
			num, err := strconv.Atoi(strings.TrimPrefix(input, "reply "))
			if err != nil {
				fmt.Println(err)
			}
			reply := replyMessage(mail[num-1])
			gomua.Send(reply)
		case input == "exit", input == "x", input == "quit", input == "q":
			exit <- true
		case strings.ContainsAny(input, "01234566789"):
			num, _ := strconv.Atoi(input)
			if num <= len(mail) && num > 0 {
				viewMessage(mail[num-1])
			}
		}

	}

	exit <- true
}

func main() {
	mail := fakeMailDir()
	viewMail(mail, os.Stdout)

	exit := make(chan bool, 1)
	go input(mail, exit)

	<-exit
	os.Exit(2)
}

// a list of static Messages for testing view
func fakeMailDir() (mail []*gomua.Message) {
	m1 := gomua.NewMessage("mr.k.frenata@gmail.com", "danweasel@gmail.com", "Sample", []string{"Message-ID: 112233445566778899"}, "Sample message content\nNew line")
	d, _ := time.ParseDuration("-2h")
	m1.SetDate(time.Now().Add(d))
	m2 := gomua.NewMessage("mr.k.frenata@gmail.com", "dan.weasel@gmail.com", "DW", []string{"Message-ID: ihu491491"}, "Breakfast?")
	m2.SetDate(time.Now())
	mail = append(mail, m1)
	mail = append(mail, m2)
	return mail
}
