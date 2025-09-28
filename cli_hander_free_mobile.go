package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func CliHandlerFreeMobile() {
	fs := flag.NewFlagSet("free_mobile", flag.ExitOnError)
	inputFile := fs.String("inputFile", "", "The input json file.")
	outputFile := fs.String("outputFile", "", "The output json file. this option required when inputFile is provided. prints to console otherwise")
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

	freeGames, err := FreeMobileGames()
	if (err != nil) {
		fmt.Println("unable to get the free mobile games")
		os.Exit(1)
	}

	if len(freeGames.ThisWeek) <= 0 {
		fmt.Println("no free games found this week")
		os.Exit(1)
	}

	if len(*inputFile) > 0 {
		// Read the original JSON file
		originalData, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Error reading:", inputFile, err)
			return
		}

		modifiedGameEntries := []map[string]interface{}{}
		if err := json.Unmarshal(originalData, &modifiedGameEntries); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		if len(freeGames.ThisWeek) == 0 {
			fmt.Println("ERROR: unexpected to see no free games this week.")
			return
		}

		// add all the free games
		for _, game := range freeGames.ThisWeek {
			modifiedEntry := map[string]interface{}{
				"epicId":        game.EpicId,
				"epicRating":    0.0,
				"epicStoreLink": game.EpicStoreLink,
				"freeDate":      game.FreeDate,
				"gameTitle":     game.GameTitle,
				"mappingSlug":   game.MappingSlug,
				"platform":      game.Platform,
				"productSlug":   game.ProductSlug,
				"sandboxId":     game.SandboxId,
				"urlSlug":       game.UrlSlug,
			}

			if len(game.SandboxId) > 0 && game.SandboxId != "TODO" {
				ratingResponse, err := RateGame(game.SandboxId)
				if err == nil {
					modifiedEntry["epicRating"] = ratingResponse.Data.RatingsPolls.GetProductResult.AverageRating
				} else {
					fmt.Println("Could not find rating for gameTitle=", game.GameTitle, "sandboxId=", game.SandboxId)
				}
			}

			modifiedGameEntries = append(modifiedGameEntries, modifiedEntry)
		}

		modifiedJSON, err := json.MarshalIndent(sortKeysInObjects(modifiedGameEntries), "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Write the modified data to the output file
		err = os.WriteFile(*outputFile, modifiedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing to:", *outputFile, err)
			return
		}

		fmt.Println("Modified data saved to ", *outputFile)
	} else {
		fmt.Println("Free This Week:")
		for _, game := range freeGames.ThisWeek {
			PrintFreeGame(game)
		}
		fmt.Println("Free Next Week:")
		for _, game := range freeGames.NextWeek {
			PrintFreeGame(game)
		}
	}
}