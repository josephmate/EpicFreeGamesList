package main

import (
	"fmt"
	"os"
)

func main() {

	// Check if enough arguments are provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: epicFreeGamesList <operation> <arguments>")
	}

	operation := os.Args[1]

	if len(operation) == 0 {
		fmt.Println("Usage: epicFreeGamesList <operation> <arguments>")
		fmt.Println("operation cannot be empty")
		return
	}

	if operation == "search" {
		CliHandlerSearch()
	} else if operation == "rate" {
		CliHandlerRating()
	} else if operation == "free" {
		CliHandlerFree()
	} else if operation == "fix_ratings" {
		CliHandlerFixRatings()
	} else {
		fmt.Println("--operation", operation, "is not recognized. only search and rate are supported")
	}
}
