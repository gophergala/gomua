package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Header Parser (hp) takes an email message as input and returns
//  the message's headers.

// For now, just take input on standard in, assuming one message
//  at a time.
// Later, allow user to specify a filename, or directory as args.

func processStdin(in io.Reader, out io.Writer) error {
	msg, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	out.Write(msg)
	return nil
}

func main() {
	if err := processStdin(os.Stdin, os.Stdout); err != nil {
		fmt.Printf("%v", err)
	}
}
