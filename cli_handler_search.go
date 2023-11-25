package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type MinimalGameEntry struct {
	FreeDate  string `json:"freeDate"`
	GameTitle string `json:"gameTitle"`
}

func CliHandlerSearch() {

	fs := flag.NewFlagSet("search", flag.ExitOnError)
	var (
		inputFile  string
		outputFile string
		gameTitle  string
	)

	fs.StringVar(&inputFile, "inputFile", "", "The input json file. --freeDate, --gameTitle cannot be used with this option")
	fs.StringVar(&outputFile, "outputFile", "", "The output json file. this option is always required")
	fs.StringVar(&gameTitle, "gameTitle", "", "The gameTitle of the free game. This option cannot be used with --inputFile")
	fs.Parse(os.Args[2:])

	if len(outputFile) == 0 && len(inputFile) > 0 {
		fmt.Println("--outputFile is required when --inputFile is used")
		return
	}
	if len(inputFile) == 0 && len(outputFile) > 0 {
		fmt.Println("--inputFile is required when --outputFile is used")
		return
	}

	if len(inputFile) > 0 && len(gameTitle) > 0 {
		fmt.Println("--inputFile cannot be used with --gameTitle")
	}

	if len(inputFile) > 0 {
		// Read the original JSON file
		originalData, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading:", inputFile, err)
			return
		}

		var gameEntries []MinimalGameEntry
		if err := json.Unmarshal(originalData, &gameEntries); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		var modifiedGameEntries []GameEntryWithSearch
		for idx, entry := range gameEntries {
			fmt.Printf("Processing: idx=%d, gameTitle=%s\n", idx, entry.GameTitle)
			modifiedEntry, err := SearchGameEntries(entry.GameTitle)
			if err != nil {
				continue
			}
			modifiedEntry.FreeDate = entry.FreeDate
			modifiedGameEntries = append(modifiedGameEntries, modifiedEntry)
		}
		// Convert the modified data to JSON
		modifiedJSON, err := json.Marshal(modifiedGameEntries)
		if err != nil {
			fmt.Println("Error converting to JSON:", err)
			return
		}

		// Write the modified data to the output file
		err = os.WriteFile(outputFile, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to:", outputFile, err)
			return
		}

		fmt.Println("Modified data saved to ", outputFile)
	} else if len(gameTitle) > 0 && len(outputFile) == 0 {
		modifiedEntry, _ := SearchGameEntries(gameTitle)
		fmt.Println("EpicId: ", modifiedEntry.EpicId)
		fmt.Println("EpicStoreLink: ", modifiedEntry.EpicStoreLink)
		fmt.Println("FreeDate: ", modifiedEntry.FreeDate)
		fmt.Println("GameTitle: ", modifiedEntry.GameTitle)
		fmt.Println("MappingSlug: ", modifiedEntry.MappingSlug)
		fmt.Println("ProductSlug: ", modifiedEntry.ProductSlug)
		fmt.Println("SandboxId: ", modifiedEntry.SandboxId)
		fmt.Println("UrlSlug: ", modifiedEntry.UrlSlug)
	} else {
		fmt.Println("--inputFile must be provided or both --gameTitle")
		return
	}
}
