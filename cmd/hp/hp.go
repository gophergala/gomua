package main

import (
	"fmt"
	"io"
	"net/mail"
	"os"
)

var exitCode = 0

// Header Parser (hp) takes an email message as input and returns
//  the message's headers.

// For now, just take input on standard in, assuming one message
//  at a time.
// Later, allow user to specify a filename, or directory as args.

func processStdin(in io.Reader, out io.Writer) error {
	msg, err := mail.ReadMessage(in)
	if err != nil {
		return err
	}

	subject := msg.Header.Get("Subject")
	fmt.Printf("%v", subject)
	return nil
}

func main() {
	if err := processStdin(os.Stdin, os.Stdout); err != nil {
		fmt.Printf("%v", err)
		exitCode = 2
	}
	os.Exit(exitCode)
}
