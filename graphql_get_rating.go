package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
	{
	  "epicId": "5d921696022d4bbe850c76be8c9bc98f",
	  "epicStoreLink": "https://store.epicgames.com/en-US/p/sunless-sea",
	  "freeDate": "2021-02-25",
	  "gameTitle": "Sunless Sea",
	  "mappingSlug": "sunless-sea",
	  "productSlug": "sunless-sea",
	  "sandboxId": "f672c20984f04f69936e6057feefe3d1",
	  "urlSlug": "rheniumgeneralaudience"
	},
*/
type GameEntryWithSearch struct {
	EpicId        string `json:"epicId"`
	EpicStoreLink string `json:"epicStoreLink"`
	FreeDate      string `json:"freeDate"`
	GameTitle     string `json:"gameTitle"`
	MappingSlug   string `json:"mappingSlug"`
	ProductSlug   string `json:"productSlug"`
	SandboxId     string `json:"sandboxId"`
	UrlSlug       string `json:"urlSlug"`
}

/*
{"data":{"RatingsPolls":{"getProductResult":{"averageRating":4.44}}},"extensions":{}}
*/
type RatingResponse struct {
	Data struct {
		RatingsPolls struct {
			GetProductResult struct {
				AverageRating float64 `json:"averageRating"`
			} `json:"getProductResult"`
		} `json:"RatingsPolls"`
	} `json:"data"`
}

func getRatingUrl(searchKey string) (RatingResponse, error) {
	format := `query%%20getProductResult($sandboxId:%%20String%%20=%%20%%22%s%%22,%%20$locale:%%20String%%20=%%20%%22US%%22)%%20{%%20RatingsPolls%%20{%%20getProductResult(sandboxId:%%20$sandboxId,%%20locale:%%20$locale)%%20{%%20averageRating%%20}%%20}%%20}`
	query := fmt.Sprintf(format, searchKey)
	// Perform the HTTP GET request to the GraphQL endpoint
	fmt.Println(query)
	resp, err := http.Get("https://graphql.epicgames.com/graphql?query=" + query)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return RatingResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return RatingResponse{}, err
	}

	fmt.Println("Processing response:", string(responseData))
	var response RatingResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		fmt.Println("Error parsing response JSON:", string(responseData))
		return RatingResponse{}, err
	}
	return response, nil
}

func RateGameEntries(gameEntries []GameEntryWithSearch, outputFile string) {
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
		response, err := getRatingUrl(searchKey)
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
	err = os.WriteFile(outputFile, modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing to:", outputFile, err)
		return
	}

	fmt.Println("Modified data saved to ", outputFile)
}
