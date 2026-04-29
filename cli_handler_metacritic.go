package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func CliHandlerMetacritic() {
	fs := flag.NewFlagSet("metacritic", flag.ExitOnError)
	inputFile := fs.String("inputFile", "", "The input json file to add Metacritic scores to")
	outputFile := fs.String("outputFile", "", "The output json file. Required when --inputFile is provided")
	gameTitle := fs.String("gameTitle", "", "The game title to search for on Metacritic")
	keepZero := fs.Bool("keepZero", false, "When set, a fetched score of 0 (TBD) will overwrite the existing score instead of preserving it")
	fs.Parse(os.Args[2:])

	if len(*outputFile) == 0 && len(*inputFile) > 0 {
		fmt.Println("--outputFile is required when --inputFile is provided")
		fs.PrintDefaults()
		os.Exit(1)
	}
	if len(*inputFile) == 0 && len(*outputFile) > 0 {
		fmt.Println("--inputFile is required when --outputFile is provided")
		fs.PrintDefaults()
		os.Exit(1)
	}

	if len(*inputFile) > 0 {
		originalData, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Error reading:", *inputFile, err)
			return
		}

		var gameEntries []GameEntryComplete
		if err := json.Unmarshal(originalData, &gameEntries); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		modifiedGameEntries := []map[string]interface{}{}

		for idx, entry := range gameEntries {
			fmt.Printf("Processing: idx=%d, gameTitle=%s\n", idx, entry.GameTitle)

			if len(entry.GameTitle) == 0 {
				fmt.Println("No gameTitle, preserving existing entry")
				modifiedGameEntries = append(modifiedGameEntries, entryToMap(entry, entry.MetacriticScore, entry.MetacriticUrl))
				continue
			}

			score, slug, err := GetMetacriticScore(entry.GameTitle)
			if err != nil {
				fmt.Println("Error getting Metacritic score, preserving existing entry:", err)
				modifiedGameEntries = append(modifiedGameEntries, entryToMap(entry, entry.MetacriticScore, entry.MetacriticUrl))
				continue
			}

			// Preserve existing score/url if the fetched score is 0, unless --keepZero is set
			metacriticUrl := ""
			if len(slug) > 0 {
				metacriticUrl = "https://www.metacritic.com/game/" + slug + "/"
			}
			if score == 0 && entry.MetacriticScore > 0 && !*keepZero {
				score = entry.MetacriticScore
				metacriticUrl = entry.MetacriticUrl
			}

			modifiedGameEntries = append(modifiedGameEntries, entryToMap(entry, score, metacriticUrl))
		}

		modifiedJSON, err := json.MarshalIndent(sortKeysInObjects(modifiedGameEntries), "", "  ")
		if err != nil {
			fmt.Println("Error converting to JSON:", err)
			return
		}

		err = os.WriteFile(*outputFile, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to:", *outputFile, err)
			return
		}

		fmt.Println("Modified data saved to", *outputFile)

	} else if len(*gameTitle) > 0 {
		score, slug, err := GetMetacriticScore(*gameTitle)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		metacriticUrl := ""
		if len(slug) > 0 {
			metacriticUrl = "https://www.metacritic.com/game/" + slug + "/"
		}
		fmt.Printf("title=%s score=%d url=%s\n", *gameTitle, score, metacriticUrl)
	} else {
		fmt.Println("Need to provide both --inputFile, --outputFile or only --gameTitle")
		fs.PrintDefaults()
		os.Exit(1)
	}
}

func entryToMap(entry GameEntryComplete, metacriticScore int, metacriticUrl string) map[string]interface{} {
	return map[string]interface{}{
		"epicId":          entry.EpicId,
		"epicRating":      entry.EpicRating,
		"epicStoreLink":   entry.EpicStoreLink,
		"freeDate":        entry.FreeDate,
		"gameTitle":       entry.GameTitle,
		"mappingSlug":     entry.MappingSlug,
		"metacriticScore": metacriticScore,
		"metacriticUrl":   metacriticUrl,
		"platform":        entry.Platform,
		"productSlug":     entry.ProductSlug,
		"sandboxId":       entry.SandboxId,
		"urlSlug":         entry.UrlSlug,
	}
}
