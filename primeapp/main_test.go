package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// need to start with Test
func Test_alpha_isPrime(t *testing.T) {

	// go test -v .

	// go test -cover .

	// go test -coverprofile=coverage.out

	// go tool cover -html=coverage.out

	// running single test
	// go test -v -run Test_checkNumbers

	// run test suits (group of tests)
	// go test -v -run Test_Alpha

	// run tests in current and in nested directories
	// go test ./...

	// table test ? vs test first
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not a prime", 8, false, "8 is not a prime number, because it is devisable by 2!"},
		{"negative", -3, false, "Negative numbers are not prime by definition!"},
		{"one", 1, false, "1 is not a prime by definition!"},
		{"zero", 0, false, "0 is not a prime by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)

		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_alpha_prompt(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to write pipe
	os.Stdout = w

	prompt() // its writing to pipe where we can save information

	// close the writer

	_ = w.Close()

	// reset os.Stdout to what it was
	os.Stdout = oldOut

	// read the output of prompt function from read pipe

	out, _ := io.ReadAll(r)

	// preform our tests, out is slice of bytes

	if string(out) != "-> " {
		t.Errorf("incorrect prompt expected -> , but got %s", string(out))
	}
}

func Test_intro(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to write pipe
	os.Stdout = w

	intro() // its writing to pipe where we can save information

	// close the writer

	_ = w.Close()

	// reset os.Stdout to what it was
	os.Stdout = oldOut

	// read the output of prompt function from read pipe

	out, _ := io.ReadAll(r)

	// preform our tests, out is slice of bytes

	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("intro text not correct, got %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: "Please enter a whole number"},
		{name: "zero", input: "0", expected: "0 is not a prime by definition!"},
		{name: "negative", input: "-3", expected: "Negative numbers are not prime by definition!"},
		{name: "one", input: "1", expected: "1 is not a prime by definition!"},
		{name: "q", input: "q", expected: ""},
		{name: "Q", input: "Q", expected: ""},
		{name: "prime", input: "7", expected: "7 is a prime number!"},
		{name: "typed", input: "three", expected: "Please enter a whole number"},
		{name: "decimal", input: "1.1", expected: "Please enter a whole number"},
	}

	for _, e := range tests {
		input := strings.NewReader(e.input)

		// create reader and pre populate it
		// can not use input, this is the way to mock it
		reader := bufio.NewScanner(input)

		res, _ := checkNumbers(reader)

		if !strings.EqualFold(e.expected, res) {
			t.Errorf("%s: expected: %s, but got: %s", e.name, e.expected, res)
		}
	}
}

func Test_readUserInput(t *testing.T) {
	// to test this function, we need a channel, and an instance of io.Reader
	doneChan := make(chan bool)

	// create a reference to bytes.Buffer
	var stdin bytes.Buffer

	// simulating user typing 1, pressing return then q and pressing return
	stdin.Write([]byte("1\nq\n"))

	go readUserInput(&stdin, doneChan)

	<-doneChan

	close(doneChan)

	// type args struct {
	// 	in       io.Reader
	// 	doneChan chan bool
	// }
	// tests := []struct {
	// 	name string
	// 	args args
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		readUserInput(tt.args.in, tt.args.doneChan)
	// 	})
	// }
}
