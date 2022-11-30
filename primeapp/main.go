package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Print a welcome message
	intro()

	// create a channel to indicate when the user wants to quit
	doneChan := make(chan bool)

	// start a gorutine to read user input and run program
	go readUserInput(os.Stdin, doneChan)

	// block until doneChan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	// say goodbye
	fmt.Println("Goodbye!")
}

func prompt() {
	fmt.Print("-> ")
}

func intro() {
	fmt.Println("Is it Prime?")
	fmt.Println("____________")
	fmt.Println("Enter a whole number, and we'll tell you if it is a Prime number or not. Enter q to quit.")
	prompt()
}

func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)

	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}

		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()

	// check if user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	// try to convert what user typed into an int
	numToCheck, err := strconv.Atoi(scanner.Text())

	if err != nil {
		return "Please enter a whole number", false
	}

	_, msg := isPrime(numToCheck)

	return msg, false
}

func isPrime(n int) (bool, string) {
	// 1 and 0 are not prime numbers
	if n == 1 || n == 0 {
		return false, fmt.Sprintf("%d is not a prime by definition!", n)
	}

	// negative numbers are not prime
	if n < 0 {
		return false, "Negative numbers are not prime by definition!"
	}

	// use the modulus operator to se if we had a prime number
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			// not a prime number
			return false, fmt.Sprintf("%d is not a prime number, because it is devisable by %d!", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number!", n)
}
