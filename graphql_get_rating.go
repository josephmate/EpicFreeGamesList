package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func RateGame(searchKey string) (RatingResponse, error) {
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
