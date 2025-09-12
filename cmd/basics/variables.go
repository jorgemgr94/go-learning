/*
Go's basic variable types:

	bool

	string

	int  int8  int16  int32  int64
	uint uint8 uint16 uint32 uint64 uintptr

	byte // alias for uint8

	rune // alias for int32
		// represents a Unicode code point

	float32 float64

	complex64 complex128
*/
package main

import "fmt"

// variables is a function that demonstrates the basic variable types in Go.
// It prints the value of the variables and the constants.
func variables() {

	fmt.Println("// == Variables ===========================================")

	// == Short declaration ==========================================
	// NOTE: in this case Go infers the type of the variable.
	messageStart := "Happy birthday! You are now"
	age := 21
	messageEnd := "years old!"

	fmt.Println(messageStart, age, messageEnd)

	// == Same line declarations ======================================
	averageOpenRate, displayMessage := .23, "is the average open rate of your messages"
	fmt.Println(averageOpenRate, displayMessage)

	// == Casting ====================================================
	accountAge := 2.6
	accountAgeInt := int(accountAge)
	fmt.Println("Your account has existed for", accountAgeInt, "years")

	// == Constants ==================================================
	const pi float64 = 3.14159
	fmt.Println("The value of pi is", pi)

	// == Computed constants ==========================================
	const secondsInMinute = 60
	const minutesInHour = 60
	const secondsInHour = secondsInMinute * minutesInHour // can compute a constant in compile time
	fmt.Println("number of seconds in an hour:", secondsInHour)

	// == Formated strings ============================================
	const name = "Saul Goodman"
	const openRate = 30.5

	msg := fmt.Sprintf("Hi %s, your open rate is %.1f percent\n", name, openRate)
	fmt.Print(msg)

	// == Conditionals ===============================================
	messageLen := 10
	maxMessageLen := 20
	fmt.Println("Trying to send a message of length:", messageLen, "and a max length of:", maxMessageLen)

	// NOTE: no parentheses are needed around the condition
	if messageLen <= maxMessageLen {
		fmt.Println("Message sent")
	} else {
		fmt.Println("Message not sent")
	}
	// NOTE: variable textLen is only available inside the if block
	// if with initial statement
	if textLen := 12; textLen <= maxMessageLen {
		fmt.Println("Message sent")
	} else {
		fmt.Println("Message not sent")
	}
}
