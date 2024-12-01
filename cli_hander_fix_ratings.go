package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func fixRatingsUsage(msg string) {
	fmt.Println("Usage: epicFreeGamesList fix_ratings <arguments>")
	fmt.Println("inputFile: The input json file. required")
	fmt.Println("outputFile: The output json file. required")
	fmt.Println(msg)
	os.Exit(1)
}

func CliHandlerFixRatings() {
	fs := flag.NewFlagSet("fix_ratings", flag.ExitOnError)
	inputFile := fs.String("inputFile", "", "The input json file. required")
	outputFile := fs.String("outputFile", "", "The output json file. required")
	fs.Parse(os.Args[2:])
	if len(*outputFile) == 0 {
		fixRatingsUsage("--outputFile is required")
	}
	if len(*inputFile) == 0 {
		fixRatingsUsage("--inputFile is required")
	}

	// Read the original JSON file
	originalData, err := os.ReadFile(*inputFile)
	if err != nil {
		fixRatingsUsage("Error reading: " + *inputFile + " " + err.Error())
	}

	var gameEntries []GameEntryComplete
	if err := json.Unmarshal(originalData, &gameEntries); err != nil {
		fixRatingsUsage("Error parsing JSON: " + err.Error())
	}

	modifiedGameEntries := []map[string]interface{}{}
	for idx, entry := range gameEntries {
		fmt.Printf("Processing: idx=%d, gameTitle=%s\n", idx, entry.GameTitle)
		
		modifiedEntry := map[string]interface{}{
			"epicId":        entry.EpicId,
			"epicRating":    entry.EpicRating,
			"epicStoreLink": entry.EpicStoreLink,
			"freeDate":      entry.FreeDate,
			"gameTitle":     entry.GameTitle,
			"mappingSlug":   entry.MappingSlug,
			"productSlug":   entry.ProductSlug,
			"sandboxId":     entry.SandboxId,
			"urlSlug":       entry.UrlSlug,
		}

		if entry.EpicRating == 0.0 {
			// rating missing. lets try to use the current information to retrieve the snapshotId that leads to the rating.
			// first lets try to use the productSlug
			mappingByPageSlug, err := GetMappingByPageSlug(entry.ProductSlug)
			if err == nil && mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId != "" {
				// use this snapshot id to get the rating
				ratingResponse, err := RateGame(mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId)
				if err == nil {
					modifiedEntry["epicId"] = mappingByPageSlug.Data.StorePageMapping.Mapping.ProductId
					modifiedEntry["epicRating"] = ratingResponse.Data.RatingsPolls.GetProductResult.AverageRating
					modifiedEntry["productSlug"] = mappingByPageSlug.Data.StorePageMapping.Mapping.PageSlug
					modifiedEntry["sandboxId"] = mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId
				}
			}

			if modifiedEntry["epicRating"] == 0.0 {
				// if we failed, lets try to use urlSlug. see README.md for the research
				mappingByPageSlug, err := GetMappingByPageSlug(entry.UrlSlug)
				if err == nil && mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId != "" {
					ratingResponse, err := RateGame(mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId)
					if err == nil {
						modifiedEntry["epicId"] = mappingByPageSlug.Data.StorePageMapping.Mapping.ProductId
						modifiedEntry["epicRating"] = ratingResponse.Data.RatingsPolls.GetProductResult.AverageRating
						modifiedEntry["productSlug"] = mappingByPageSlug.Data.StorePageMapping.Mapping.PageSlug
						modifiedEntry["sandboxId"] = mappingByPageSlug.Data.StorePageMapping.Mapping.SandboxId
					}
				}
			}
		}

		modifiedGameEntries = append(modifiedGameEntries, modifiedEntry)
	}

	// Convert the modified data to JSON
		modifiedJSON, err := json.MarshalIndent(modifiedGameEntries, "", "  ")
	if err != nil {
		fixRatingsUsage("Error converting to JSON: " + err.Error())
		return
	}

	// Write the modified data to the output file
	err = os.WriteFile(*outputFile, modifiedJSON, 0644)
	if err != nil {
		fixRatingsUsage("Error writing to: " + *outputFile + " " + err.Error())
		return
	}

	fmt.Println("Modified data saved to ", *outputFile)
}
