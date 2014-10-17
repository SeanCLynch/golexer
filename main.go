package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Regex for identifying input tokens

var oct, _ = regexp.Compile("[0-7]")
var dec, _ = regexp.Compile("[8-9]")
var hex, _ = regexp.Compile("[aAc-fC-F]")
var sign, _ = regexp.Compile("[-+]")
var octnum, _ = regexp.Compile("[bB]")
var hexnum, _ = regexp.Compile("[hH]")

// State is used to, well, keep track of the state

type State int

// Token is what is emitted by analyse()

type Token struct {
	input     string
	length    int
	error     string
	value     int64
	base      int
	sign      bool
	accepting bool
}

// Helper function for printing tests/errors

func ttype(i int) string {
	switch {
	case i == 1:
		return "Octal"
	case i == 2:
		return "Decimal"
	case i == 3:
		return "Hex"
	case i == 4:
		return "Sign"
	case i == 5:
		return "Hex B"
	case i == 6:
		return "Oct H"
	case i == 7:
		return "END"
	case i == 8:
		return "INVALID"
	}
	return "HUH?"
}

// nextState calculates the next state given the previous state
// and the type of the current character. The token is used
// to create special actions for each transition.

func nextState(s *State, out *Token, c string, t int) State {
	if strings.EqualFold(os.Args[1], "tests") {
		fmt.Printf("State:%d -- %s (%s)", int(*s), c, ttype(t))
	}
	switch {
	case int(*s) == 1:

		if t == 1 { //oct
			*s = State(2)
			out.base = 10
		} else if t == 2 { //dec
			*s = State(3)
			out.base = 10
		} else if t == 3 { //hex
			*s = State(4)
			out.base = 16
		} else if t == 4 { //sign
			*s = State(5)
			out.base = 10
		} else if t == 5 { //octnum b
			*s = State(4)
			out.base = 16
		} else {
			*s = State(0) // Error State
			out.error = "Invalid starting character."
		}

	case int(*s) == 2:

		if t == 1 { //oct
			*s = State(2)
			out.base = 10
		} else if t == 2 { //dec
			*s = State(3)
			out.base = 10
		} else if t == 3 { //hex
			*s = State(4)
			out.base = 16
		} else if t == 5 { //octnum b
			*s = State(7)
			out.base = 8
		} else if t == 6 { //hexnum h
			*s = State(6)
			out.base = 16
		} else if t == 7 {
			*s = State(8) // Accepting State
			out.accepting = true
		} else {
			*s = State(0) // Error State
			out.error = "Error in State 2"
		}

	case int(*s) == 3:

		if t == 1 { //oct
			*s = State(3)
			out.base = 10
		} else if t == 2 { //dec
			*s = State(3)
			out.base = 10
		} else if t == 3 { //hex
			if out.sign == true {
				*s = State(0)
				out.base = 10
				out.error = "Hex value after sign"
			} else {
				*s = State(4)
				out.base = 16
			}
		} else if t == 5 { //octnum b
			if out.sign == true {
				*s = State(0)
				out.base = 10
				out.error = "Hex value after sign"
			} else {
				*s = State(4)
				out.base = 16
			}
		} else if t == 6 { //hexnum h
			if out.sign == true {
				*s = State(0)
				out.base = 10
				out.error = "Hex end line after sign"
			} else {
				*s = State(6)
				out.base = 16
			}
		} else if t == 7 {
			*s = State(8) // Accepting State
			out.accepting = true
		} else {
			*s = State(0) // Error State
			out.error = "Error in State 3"
		}

	case int(*s) == 4:

		if t == 1 { //oct
			*s = State(4)
		} else if t == 2 { //dec
			*s = State(4)
		} else if t == 3 { //hex
			*s = State(4)
		} else if t == 5 { //octnum b
			*s = State(4)
		} else if t == 6 { //hexnum h
			*s = State(6)
			out.base = 16
		} else {
			*s = State(0) // Error State
			out.error = "Error in State 4"
		}

	case int(*s) == 5:

		out.sign = true
		if t == 1 { //oct
			*s = State(3)
			out.base = 10
		} else if t == 2 { //dec
			*s = State(3)
			out.base = 10
		} else {
			*s = State(0) // Error State
			out.error = "Error in State 5"
		}

	case int(*s) == 6:

		if t == 7 {
			*s = State(8) // Accepting State
			out.accepting = true
		} else {
			*s = State(0) // Error State
			out.error = "non-end line character after [hH]"
		}

	case int(*s) == 7:

		if t == 1 { //oct
			*s = State(4)
			out.base = 16
		} else if t == 2 { //dec
			*s = State(4)
			out.base = 16
		} else if t == 3 { //hex
			*s = State(4)
			out.base = 16
		} else if t == 5 { //octnum b
			*s = State(4)
			out.base = 16
		} else if t == 6 { //hexnum h
			*s = State(6)
			out.base = 16
		} else if t == 7 {
			*s = State(8) // Accepting State
			out.accepting = true
		} else {
			*s = State(0) // Error State
			out.error = "Error in State 7"
		}

	}
	if strings.EqualFold(os.Args[1], "tests") {
		fmt.Printf(" -> State: %d \n", int(*s))
	}
	return *s
}

// The analyse function takes an input string and creates a token
// to return to the user. It iterates through each character, then
// calls nextState() for that character and the current state.
// At the end of the input, it calculates the actual value, and
// determines if that value is too large.

func analyze(str string) Token {
	currState := State(1)
	currTok := 0
	retTok := Token{}
	retTok.input = str
	retTok.length = len(str)

	for i := 0; i < len(str); i++ {
		s := string(str[i])
		if oct.MatchString(s) {
			currTok = 1
		} else if dec.MatchString(s) {
			currTok = 2
		} else if hex.MatchString(s) {
			currTok = 3
		} else if sign.MatchString(s) {
			currTok = 4
		} else if octnum.MatchString(s) {
			currTok = 5
		} else if hexnum.MatchString(s) {
			currTok = 6
		} else {
			currTok = 8
		}
		currState = nextState(&currState, &retTok, s, currTok)
		if i == (len(str) - 1) {
			currTok = 7
			currState = nextState(&currState, &retTok, s, currTok)
		}
	}

	if retTok.base == 8 {
		tempval := retTok.input[:len(retTok.input)-1]
		value, _ := strconv.ParseInt(tempval, retTok.base, 64)
		if value > 4294967295 {
			retTok.accepting = false
			retTok.error = "Too large #"
		}
		retTok.value = value
	} else if retTok.base == 10 {
		tempval := "0"
		if retTok.sign {
			tempval = retTok.input[1:]
		} else {
			tempval = retTok.input
		}
		value, _ := strconv.ParseInt(tempval, retTok.base, 64)
		retTok.value = value
	} else if retTok.base == 16 {
		tempval := retTok.input[:len(retTok.input)-1]
		value, _ := strconv.ParseInt(tempval, retTok.base, 64)
		if value > 4294967295 {
			retTok.accepting = false
			retTok.error = "Too large #"
		}
		retTok.value = value
	}

	// fmt.Printf("Analyzed %s [state: %d ] \n", str, int(currState))
	fmt.Printf("EMIT TOKEN: [accept:%t] [const:%s] [base:%d] [value:%d] [error:%s] \n\n", retTok.accepting, retTok.input, retTok.base, retTok.value, retTok.error)
	return retTok
}

// The main function simply checks if the input is the keyword
// tests - and executes the test cases if it is.

func main() {
	str := os.Args[1]

	if !strings.EqualFold(os.Args[1], "tests") {
		analyze(str)
	}

	if strings.EqualFold(os.Args[1], "tests") {
		fmt.Println("VALIDS: ----------------------------------")
		analyze("5")
		analyze("9")
		analyze("23h")
		analyze("+4")
		analyze("-8")
		analyze("dh")
		analyze("bh")
		analyze("4b")
		analyze("3bh")
		analyze("3ah")
		analyze("4b2h")
		analyze("1b8h")
		analyze("2bch")
		analyze("3bbh")
		analyze("1859h")
		analyze("8ach")
		analyze("8b4h")
		analyze("8c9h")
		analyze("8bbh")
		fmt.Println("INVALIDS: ----------------------------------")
		analyze("-3b")
		analyze("2c")
		analyze("+3h")
		analyze("-2ch")
	}

}
