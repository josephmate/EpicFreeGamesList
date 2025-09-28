package main

import (
	"fmt"
	"os"
	"strconv"
)


func mainUsage(msg string) {
		fmt.Println("Usage: epicFreeGamesList <operation> <arguments>")
		fmt.Println("  search: look for a game by string")
		fmt.Println("  rate: find the rating of game given some identifier")
		fmt.Println("  free: update the free game list with the pc games")
		fmt.Println("  free_mobile: update the free game list with the mobile games")
		fmt.Println("  fix_ratings: fix any broken ratings that can be fixed")
		fmt.Println("  mobile_discover_page: print the games from the mobile page")

		fmt.Println(msg)
		os.Exit(1)
}


func main() {

	// Check if enough arguments are provided
	if len(os.Args) < 2 {
		mainUsage("not enough arguments. only had " + strconv.Itoa(len(os.Args)))
	}

	operation := os.Args[1]

	if len(operation) == 0 {
		mainUsage("operation cannot be empty")
	}

	if operation == "search" {
		CliHandlerSearch()
	} else if operation == "rate" {
		CliHandlerRating()
	} else if operation == "free" {
		CliHandlerFree()
	} else if operation == "free_mobile" {
		CliHandlerFreeMobile()
	} else if operation == "fix_ratings" {
		CliHandlerFixRatings()
	} else if operation == "mobile_discover_page" {
		CliHandlerMobileDiscoverPage()
	} else {
		mainUsage("--operation " + operation + " is not recognized.")
	}
}
