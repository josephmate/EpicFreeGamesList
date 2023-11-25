package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func CliHandlerRating() {
	fs := flag.NewFlagSet("rating", flag.ExitOnError)
	inputFile := fs.String("inputFile", "", "The input json file. --freeDate, --gameTitle cannot be used with this option")
	outputFile := fs.String("outputFile", "", "The output json file. this option is always required")
	searchKey := fs.String("searchKey", "", "The searchKey used to search for a rating.")
	fs.Parse(os.Args[2:])
	if len(*outputFile) == 0 && len(*inputFile) > 0 {
		fmt.Println("--outputFile is required when --inputFile is provided")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(*inputFile) == 0 && len(*outputFile) > 0 {
		fmt.Println("--inputFile is required when --outputFile is provided")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(*inputFile) > 0 {
		// Read the original JSON file
		originalData, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Error reading:", *inputFile, err)
			return
		}

		var gameEntries []GameEntryWithSearch
		if err := json.Unmarshal(originalData, &gameEntries); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		modifiedGameEntries := []map[string]interface{}{}

		for idx, entry := range gameEntries {
			fmt.Printf("Processing: idx=%d, gameTitle=%s\n", idx, entry.GameTitle)

			var searchKey string
			if len(entry.SandboxId) > 0 {
				searchKey = entry.SandboxId
			} else if len(entry.MappingSlug) > 0 {
				searchKey = entry.MappingSlug
			} else if len(entry.ProductSlug) > 0 {
				searchKey = entry.ProductSlug
			} else if len(entry.UrlSlug) > 0 {
				searchKey = entry.UrlSlug
			} else if len(entry.EpicId) > 0 {
				searchKey = entry.EpicId
			} else {
				jsonStr, _ := json.Marshal(entry)
				fmt.Println("No useful property provided on :", jsonStr)
				continue
			}
			response, err := RateGame(searchKey)
			if err != nil {
				fmt.Println("Error getting rating:", err)
				continue
			}

			// Create the modified entry
			modifiedEntry := map[string]interface{}{
				"epicId":        entry.EpicId,
				"epicRating":    response.Data.RatingsPolls.GetProductResult.AverageRating,
				"epicStoreLink": entry.EpicStoreLink,
				"freeDate":      entry.FreeDate,
				"gameTitle":     entry.GameTitle,
				"mappingSlug":   entry.MappingSlug,
				"productSlug":   entry.ProductSlug,
				"sandboxId":     entry.SandboxId,
				"urlSlug":       entry.UrlSlug,
			}

			modifiedGameEntries = append(modifiedGameEntries, modifiedEntry)
		}

		// Convert the modified data to JSON
		modifiedJSON, err := json.Marshal(modifiedGameEntries)
		if err != nil {
			fmt.Println("Error converting to JSON:", err)
			return
		}

		// Write the modified data to the output file
		err = os.WriteFile(*outputFile, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to:", outputFile, err)
			return
		}

		fmt.Println("Modified data saved to ", outputFile)

	} else if len(*searchKey) > 0 {
		rating, _ := RateGame(*searchKey)
		fmt.Printf("%+v\n", rating)
	} else {
		fmt.Println("Need to provide both --inputFile, --outputFile or only --searchKey")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
