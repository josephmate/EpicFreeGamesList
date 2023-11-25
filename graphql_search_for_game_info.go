package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type SearchResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []struct {
					UrlSlug     string `json:"urlSlug"`
					ProductSlug string `json:"productSlug"`
					Id          string `json:"id"`
					Namespace   string `json:"namespace"`
					CatalogNs   struct {
						Mappings []struct {
							PageSlug string `json:"pageSlug"`
						} `json:"mappings"`
					} `json:"catalogNs"`
				} `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

func SearchGameEntries(gameTitle string) (GameEntryWithSearch, error) {
	// Create the GraphQL query with entry.gameTitle
	query := fmt.Sprintf(`query%%20searchStoreQuery($allowCountries:%%20String,%%20$category:%%20String,%%20$locale:%%20String,%%20$namespace:%%20String,%%20$itemNs:%%20String,%%20$sortBy:%%20String,%%20$sortDir:%%20String,%%20$start:%%20Int,%%20$tag:%%20String,%%20$releaseDate:%%20String,%%20$priceRange:%%20String,%%20$freeGame:%%20Boolean,%%20$onSale:%%20Boolean,%%20$effectiveDate:%%20String)%%20{%%20Catalog%%20{%%20searchStore(%%20allowCountries:%%20$allowCountries%%20category:%%20$category%%20count:%%201%%20country:%%20"US"%%20keywords:%%20%%22%s%%22%%20locale:%%20$locale%%20namespace:%%20$namespace%%20itemNs:%%20$itemNs%%20sortBy:%%20$sortBy%%20sortDir:%%20$sortDir%%20releaseDate:%%20$releaseDate%%20start:%%20$start%%20tag:%%20$tag%%20priceRange:%%20$priceRange%%20freeGame:%%20$freeGame%%20onSale:%%20$onSale%%20effectiveDate:%%20$effectiveDate%%20)%%20{%%20elements%%20{%%20title%%20id%%20namespace%%20description%%20effectiveDate%%20productSlug%%20urlSlug%%20url%%20tags%%20{%%20id%%20}%%20items%%20{%%20id%%20namespace%%20}%%20customAttributes%%20{%%20key%%20value%%20}%%20categories%%20{%%20path%%20}%%20catalogNs%%20{%%20mappings(pageType:%%20"productHome")%%20{%%20pageSlug%%20pageType%%20}%%20}%%20offerMappings%%20{%%20pageSlug%%20pageType%%20}%%20}%%20}%%20}%%20}`,
		url.QueryEscape(gameTitle))

	// Perform the HTTP GET request to the GraphQL endpoint
	resp, err := http.Get("https://graphql.epicgames.com/graphql?query=" + query)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return GameEntryWithSearch{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return GameEntryWithSearch{}, err
	}

	var response SearchResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		fmt.Println("Error parsing response JSON:", string(responseData))
		return GameEntryWithSearch{}, err
	}

	// Create the modified entry
	modifiedEntry := GameEntryWithSearch{
		EpicStoreLink: "TODO",
		FreeDate:      "TODO",
		GameTitle:     gameTitle,
		EpicId:        "TODO",
		UrlSlug:       "TODO",
		ProductSlug:   "TODO",
		MappingSlug:   "TODO",
		SandboxId:     "TODO",
	}

	fmt.Println("Processing response:", string(responseData))
	tmpJsonToBytes, _ := json.Marshal(response)
	fmt.Println("Converted to json back to string", string(tmpJsonToBytes))
	if len(response.Data.Catalog.SearchStore.Elements) > 0 {
		element := response.Data.Catalog.SearchStore.Elements[0]
		modifiedEntry.EpicId = element.Id
		mapping := element.CatalogNs.Mappings
		if len(mapping) > 0 {
			modifiedEntry.EpicStoreLink = "https://store.epicgames.com/en-US/p/" + mapping[0].PageSlug
			modifiedEntry.MappingSlug = mapping[0].PageSlug
		} else if strings.TrimSpace(element.ProductSlug) != "" {
			modifiedEntry.EpicStoreLink = "https://store.epicgames.com/en-US/p/" + strings.TrimSuffix(element.ProductSlug, "/home")
		} else if strings.TrimSpace(element.UrlSlug) != "" {
			// urlSlug is the last resort since some times is had made up data.
			modifiedEntry.EpicStoreLink = "https://store.epicgames.com/en-US/p/" + element.UrlSlug
		} else {
			fmt.Println("Did not have pageSlug at all")
		}

		modifiedEntry.UrlSlug = element.UrlSlug
		modifiedEntry.ProductSlug = element.ProductSlug
		modifiedEntry.SandboxId = response.Data.Catalog.SearchStore.Elements[0].Namespace
	} else {
		fmt.Println("Did not have Elements")
	}

	return modifiedEntry, nil
}
