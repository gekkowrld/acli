package errors

import (
	"fmt"
	"os"
)

type errorType struct {
	Error     string
	ErrorCode int
}

// Initialize and populate the error map
var errMap = map[int]errorType{
	1: {"Unkown Error Occured", 1},
	2: {"Permission denied", 2},
	3: {"Not a git repository", 3},
}

// Display errors and exit after that
func ExitError(errorCode int) {
	if err, ok := errMap[errorCode]; ok {
		fmt.Println(err.Error)
		os.Exit(errorCode)
	} else {
		fmt.Println("Unknown Error Occured")
		os.Exit(1)
	}
}

// Display errors and do nothing
func DispError(errorCode int) {
	if err, ok := errMap[errorCode]; ok {
		fmt.Println(err.Error)
	} else {
		fmt.Println("Unknown Error Occured")
	}
}
