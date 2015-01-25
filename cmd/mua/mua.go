package main

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func scanMailDir(dir string) (mail []*Message) {

}

// takes a slice of Messages and prints a numbered list of summaries
func viewMail(mail []*Message) {

}

// prints a single mail message to the screen
func viewMessage(msg *Message) {

}

// (to be invoked from viewMessage with a keypress, most likely)
// create a new message:
// put old's Message-ID into reply's In-Reply-To and References headers
// put old's content as quoted material in reply
// put old's From as the reply's To
// put old's Subject in reply's subject as "RE: Subject" (but allow user to edit? -- later)
// prompt user for content
// send reply
func replyMessage(old *Message) (reply *Message) {

	return reply
}

func main() {
	mail := fakeMailDir()
	viewMail(mail)
}

// a list of static Messages for testing view
func fakeMailDir() (mail []*Message) {

	return mail
}
