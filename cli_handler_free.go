package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "sort"
)

func PrintFreeGame(game FreeGameEntry) {
    fmt.Println("EpicId: ", game.EpicId)
    fmt.Println("EpicStoreLink: ", game.EpicStoreLink)
    fmt.Println("FreeDate: ", game.FreeDate)
    fmt.Println("GameTitle: ", game.GameTitle)
    fmt.Println("MappingSlug: ", game.MappingSlug)
    fmt.Println("Platform: ", game.Platform)
    fmt.Println("ProductSlug: ", game.ProductSlug)
    fmt.Println("SandboxId: ", game.SandboxId)
    fmt.Println("UrlSlug: ", game.UrlSlug)
    fmt.Println("------------------------------------")
}

func sortKeysInObjects(input []map[string]interface{}) []map[string]interface{} {
    for _, obj := range input {
        sortedKeys := make([]string, 0, len(obj))
        for key := range obj {
            sortedKeys = append(sortedKeys, key)
        }
        sort.Strings(sortedKeys)

        sortedData := make(map[string]interface{})
        for _, key := range sortedKeys {
            sortedData[key] = obj[key]
        }

        obj = sortedData
    }
    return input
}

func CliHandlerFree() {
    fs := flag.NewFlagSet("free", flag.ExitOnError)
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

    freeGames, _ := GetFreeGames()

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
