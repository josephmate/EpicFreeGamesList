package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

/*
EpicId:  *freeDate,
EpicStoreLink:  *freeDate,
FreeDate:  *freeDate,
GameTitle: *gameTitle,
MappingSlug:  *freeDate,
ProductSlug:  *freeDate,
SandboxId:  *freeDate,
UrlSlug:  *freeDate,
*/
var (
	operation     = flag.String("op", "", "The input json file. --freeDate, --gameTitle cannot be used with this option")
	inputFile     = flag.String("inputFile", "", "The input json file. --freeDate, --gameTitle cannot be used with this option")
	outputFile    = flag.String("outputFile", "", "The output json file. this option is always required")
	epicId        = flag.String("epicId", "", "The id from epic games. This option cannot be used with --inputFile")
	epicStoreLink = flag.String("epicStoreLink", "", "The url that game can be found on at epic. This option cannot be used with --inputFile")
	freeDate      = flag.String("freeDate", "", "The date the game was free on. This option cannot be used with --inputFile")
	gameTitle     = flag.String("gameTitle", "", "The gameTitle of the free game. This option cannot be used with --inputFile")
	mappingSlug   = flag.String("mappingSlug", "", "This option cannot be used with --inputFile")
	productSlug   = flag.String("productSlug", "", "This option cannot be used with --inputFile")
	sandboxId     = flag.String("sandboxId", "", "This option cannot be used with --inputFile")
	urlSlug       = flag.String("urlSlug", "", "This option cannot be used with --inputFile")
)

func main() {

	flag.Parse()
	if len(*operation) == 0 {
		fmt.Println("--operation is always required")
		return
	}
	if len(*outputFile) == 0 {
		fmt.Println("--outputFile is always required")
		return
	}

	if *operation == "search" {
		if len(*inputFile) > 0 && len(*freeDate) > 0 {
			fmt.Println("--inputFile cannot be used with --freeDate")
			return
		}
		if len(*inputFile) > 0 && len(*gameTitle) > 0 {
			fmt.Println("--inputFile cannot be used with --gameTitle")
		}

		if len(*inputFile) > 0 {
			// Read the original JSON file
			originalData, err := os.ReadFile(*inputFile)
			if err != nil {
				fmt.Println("Error reading:", *inputFile, err)
				return
			}

			var gameEntries []MinimalGameEntry
			if err := json.Unmarshal(originalData, &gameEntries); err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}
			SearchGameEntries(gameEntries, *outputFile)
		} else if len(*gameTitle) > 0 && len(*freeDate) > 0 {
			var gameEntries []MinimalGameEntry
			gameEntries = append(gameEntries, MinimalGameEntry{
				FreeDate:  *freeDate,
				GameTitle: *gameTitle,
			})
			SearchGameEntries(gameEntries, *outputFile)
		} else if len(*gameTitle) > 0 {
			fmt.Println("--freeDate must be used with --gameTitle")
			return
		} else if len(*freeDate) > 0 {
			fmt.Println("--gameTitle must be used with --freeDate")
			return
		} else {
			fmt.Println("--inputFile must be provided or both --gameTitle and --freeDate")
			return
		}
	} else if *operation == "rate" {
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
			RateGameEntries(gameEntries, *outputFile)
		} else {
			var gameEntries []GameEntryWithSearch
			gameEntries = append(gameEntries, GameEntryWithSearch{
				EpicId:        *epicId,
				EpicStoreLink: *epicStoreLink,
				FreeDate:      *freeDate,
				GameTitle:     *gameTitle,
				MappingSlug:   *mappingSlug,
				ProductSlug:   *productSlug,
				SandboxId:     *sandboxId,
				UrlSlug:       *urlSlug,
			})
			RateGameEntries(gameEntries, *outputFile)
		}
	} else {
		fmt.Println("--operation", *operation, "is not recognized. only search and rate are supported")
	}
}
