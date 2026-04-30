package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func CliHandlerSteamDB() {
	fs := flag.NewFlagSet("steamdb", flag.ExitOnError)
	inputFile := fs.String("inputFile", "", "The input json file to add SteamDB ratings to")
	outputFile := fs.String("outputFile", "", "The output json file. Required when --inputFile is provided")
	gameTitle := fs.String("gameTitle", "", "The game title to search for on Steam")
	keepZero := fs.Bool("keepZero", false, "When set, a fetched rating of 0 will overwrite the existing rating instead of preserving it")
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
				modifiedGameEntries = append(modifiedGameEntries, entryToMapSteamDB(entry, entry.SteamDBRating, entry.SteamDBUrl, entry.SteamUrl))
				continue
			}

			rating, appID, err := GetSteamDBRating(entry.GameTitle)
			if err != nil {
				fmt.Println("Error getting SteamDB rating, preserving existing entry:", err)
				modifiedGameEntries = append(modifiedGameEntries, entryToMapSteamDB(entry, entry.SteamDBRating, entry.SteamDBUrl, entry.SteamUrl))
				continue
			}

			steamDBUrl := ""
			steamUrl := ""
			if appID > 0 {
				steamDBUrl = fmt.Sprintf("https://steamdb.info/app/%d/", appID)
				steamUrl = fmt.Sprintf("https://store.steampowered.com/app/%d/", appID)
			}

			// Preserve existing rating/urls if the fetched rating is 0, unless --keepZero is set.
			if rating == 0 && entry.SteamDBRating > 0 && !*keepZero {
				rating = entry.SteamDBRating
				steamDBUrl = entry.SteamDBUrl
				steamUrl = entry.SteamUrl
			}

			modifiedGameEntries = append(modifiedGameEntries, entryToMapSteamDB(entry, rating, steamDBUrl, steamUrl))
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
		rating, appID, err := GetSteamDBRating(*gameTitle)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		steamDBUrl := ""
		steamUrl := ""
		if appID > 0 {
			steamDBUrl = fmt.Sprintf("https://steamdb.info/app/%d/", appID)
			steamUrl = fmt.Sprintf("https://store.steampowered.com/app/%d/", appID)
		}
		fmt.Printf("title=%s rating=%.2f steamDBUrl=%s steamUrl=%s\n", *gameTitle, rating, steamDBUrl, steamUrl)
	} else {
		fmt.Println("Need to provide both --inputFile, --outputFile or only --gameTitle")
		fs.PrintDefaults()
		os.Exit(1)
	}
}

func entryToMapSteamDB(entry GameEntryComplete, steamDBRating float64, steamDBUrl string, steamUrl string) map[string]interface{} {
	return map[string]interface{}{
		"epicId":          entry.EpicId,
		"epicRating":      entry.EpicRating,
		"epicStoreLink":   entry.EpicStoreLink,
		"freeDate":        entry.FreeDate,
		"gameTitle":       entry.GameTitle,
		"mappingSlug":     entry.MappingSlug,
		"metacriticScore": entry.MetacriticScore,
		"metacriticUrl":   entry.MetacriticUrl,
		"platform":        entry.Platform,
		"productSlug":     entry.ProductSlug,
		"sandboxId":       entry.SandboxId,
		"steamDBRating":   steamDBRating,
		"steamDBUrl":      steamDBUrl,
		"steamUrl":        steamUrl,
		"urlSlug":         entry.UrlSlug,
	}
}
