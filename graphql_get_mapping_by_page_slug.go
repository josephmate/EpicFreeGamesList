package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
{
  "data": {
    "StorePageMapping": {
      "mapping": {
        "pageSlug": "payday-2-c66369",
        "pageType": "productHome",
        "sandboxId": "3b661fd6a9724ac7b6ac6d10d0572511",
        "productId": "14eb3477a6084940b49de5aa73c60f98",
        "createdDate": "2023-06-07T08:05:53.761Z",
        "updatedDate": "2023-11-27T20:21:20.490Z",
        "mappings": {
          "cmsSlug": null,
          "pageId": null
        }
      }
    }
  },
  "extensions": {}
}
*/
type StorePageMapping struct {
	Data struct {
		StorePageMapping struct {
			Mapping struct {
				PageSlug string `json:pageSlug`
				SandboxId string `json:sandboxId`
				ProductId string `json:productId`
			} `json:"mapping"`
		} `json:"StorePageMapping"`
	} `json:"data"`
}

/*
query%20getMappingByPageSlug%28%24pageSlug%3A%20String%21%20%3D%20%22payday-2-c66369%22%2C%20%24sandboxId%3A%20String%29%20%7B%0A%20%20StorePageMapping%20%7B%0A%20%20%20%20mapping%28pageSlug%3A%20%24pageSlug%2C%20sandboxId%3A%20%24sandboxId%29%20%7B%0A%20%20%20%20%20%20pageSlug%0A%20%20%20%20%20%20pageType%0A%20%20%20%20%20%20sandboxId%0A%20%20%20%20%20%20productId%0A%20%20%20%20%20%20createdDate%0A%20%20%20%20%20%20updatedDate%0A%20%20%20%20%20%20mappings%20%7B%0A%20%20%20%20%20%20%20%20cmsSlug%0A%20%20%20%20%20%20%20%20pageId%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D
*/
func GetMappingByPageSlug(searchKey string) (StorePageMapping, error) {
	format := `query%%20getMappingByPageSlug%%28%%24pageSlug%%3A%%20String%%21%%20%%3D%%20%%22%s%%22%%2C%%20%%24sandboxId%%3A%%20String%%29%%20%%7B%%0A%%20%%20StorePageMapping%%20%%7B%%0A%%20%%20%%20%%20mapping%%28pageSlug%%3A%%20%%24pageSlug%%2C%%20sandboxId%%3A%%20%%24sandboxId%%29%%20%%7B%%0A%%20%%20%%20%%20%%20%%20pageSlug%%0A%%20%%20%%20%%20%%20%%20pageType%%0A%%20%%20%%20%%20%%20%%20sandboxId%%0A%%20%%20%%20%%20%%20%%20productId%%0A%%20%%20%%20%%20%%20%%20createdDate%%0A%%20%%20%%20%%20%%20%%20updatedDate%%0A%%20%%20%%20%%20%%20%%20mappings%%20%%7B%%0A%%20%%20%%20%%20%%20%%20%%20%%20cmsSlug%%0A%%20%%20%%20%%20%%20%%20%%20%%20pageId%%0A%%20%%20%%20%%20%%20%%20%%7D%%0A%%20%%20%%20%%20%%7D%%0A%%20%%20%%7D%%0A%%7D`
	query := fmt.Sprintf(format, searchKey)
	// Perform the HTTP GET request to the GraphQL endpoint
	fmt.Println(query)
	resp, err := http.Get("https://graphql.epicgames.com/graphql?query=" + query)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return StorePageMapping{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return StorePageMapping{}, err
	}

	fmt.Println("Processing response:", string(responseData))
	var response StorePageMapping
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		fmt.Println("Error parsing response JSON:", string(responseData))
		return StorePageMapping{}, err
	}
	return response, nil
}
