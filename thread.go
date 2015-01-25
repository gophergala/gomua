package gomua

import (
	"fmt"
)

// Take a slice of gomua.Messages and sort them into threads
// A thread starts as a linked-list of gomua.Messages

// Display a threaded conversation
func (m *ThreadedMessage) String() string {
	var output string = fmt.Sprintf("From: %v\n", m.Header.Get("From")) +
		fmt.Sprintf("To: %v\n", m.Header.Get("To")) +
		fmt.Sprintf("Date: %v\n", m.Header.Get("Date")) +
		fmt.Sprintf("Subject: %v\n", m.Header.Get("Subject")) +
		fmt.Sprintf("\n%s\n", m.Content)

	return output
}

// Summary gets a summarized subject for this thread (i.e. remove Re: re: Re:)
func (m *ThreadedMessage) Summary() string {
	subject := m.Header.Get("Subject")
	from := m.Header.Get("From")
	return fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))
}

func Thread(msgs []Mail) []Mail {
	// take a slice of mails
	// thread them!
	for i, m := range msgs {
		// check the subject -- if we've seen it, put it in the list
		fmt.Printf("test: %d %v", i, m.Summary())
		return nil
		//fmt.Printf("%d. %s\n", i+1, m.Summary())
	}
	fmt.Printf("Hello\n\n")
	return msgs
}
