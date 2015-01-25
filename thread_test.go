package gomua

import (
	"fmt"
	"io"
	"testing"
)

// Test threading

// Get messages

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func scanMailDir(dir string) (msgs []Mail) {
	mails := Scan(dir)

	// Embed mail.Message inside gomua.Message
	for _, m := range mails {
		msg := ReadMessage(m)
		msgs = append(msgs, msg)
	}

	return msgs
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMailList(msgs []Mail, w io.Writer) {
	for i, m := range msgs {
		fmt.Fprintf(w, "%d. %s\n", i+1, m.Summary())
	}
}

func TestGetMessages(t *testing.T) {
	fmt.Print("Hello")
	const dir string = "./cmd/mua/testmaildir"
	msgs := scanMailDir(dir)
	msgs = Thread(msgs)
}
