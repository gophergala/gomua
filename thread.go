package gomua

import (
	"fmt"
)

// MessageThread is a linked list of Messages
type MessageThread struct {
	head *ThreadNode
}

type ThreadNode struct {
	msg  *Message
	next *ThreadNode
}

var threads = make(map[string]*MessageThread)

// Take a slice of gomua.Messages and sort them into threads
// A thread starts as a linked-list of gomua.Messages

// Display a threaded conversation
func (t *MessageThread) String() string {
	var output string = fmt.Sprintf("From: %v\n", t.head.msg.Header.Get("From")) +
		fmt.Sprintf("To: %v\n", t.head.msg.Header.Get("To")) +
		fmt.Sprintf("Date: %v\n", t.head.msg.Header.Get("Date")) +
		fmt.Sprintf("Subject: %v\n", t.head.msg.Header.Get("Subject")) +
		fmt.Sprintf("\n%s\n", t.head.msg.Content)

	return output
}

// Summary gets a summarized subject for this thread (i.e. remove Re: re: Re:)
func (t *MessageThread) Summary() string {
	subject := t.head.msg.Header.Get("Subject")
	from := t.head.msg.Header.Get("From")
	return fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))
}

func (t *MessageThread) appendNode(n *ThreadNode) {
	node := t.head
	for node.next != nil {
		node = node.next
	}
	node.next = n
}

func Thread(msgs []Mail) []Mail {
	// take a slice of mails
	// make a hash table keyed off of subject for now

	//thread := new(MessageThread)
	// thread them!
	for i, m := range msgs {
		// check the subject -- if we've seen it, put it in the list
		node := new(ThreadNode)
		node.msg = m.(*Message)
		if threads[m.Summary()] == nil {
			fmt.Printf("New thread (%d): '%v'\n", i, m.Summary())
			thread := new(MessageThread)
			thread.head = m
			threads[m.Summary()] = thread
		} else {
			threads[m.Summary()] =
				fmt.Printf("Add to thread (%d): '%v'\n", i, m.Summary())
			thread.appendNode(node)
		}
	}
	return msgs
}
