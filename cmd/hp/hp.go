package main

import (
	"io"
)

// Header Parser (hp) takes an email message as input and returns
//  the message's headers.

// For now, just take input on standard in, assuming one message
//  at a time.
// Later, allow user to specify a filename, or directory as args.

func processStdin(in io.Reader, out io.Writer) error {
	return nil
}

func main() {
}
