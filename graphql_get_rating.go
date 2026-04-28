package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	http "github.com/bogdanfinn/fhttp"
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
	variables := fmt.Sprintf(`{"sandboxId":"%s","locale":"en-US"}`, searchKey)
	extensions := `{"persistedQuery":{"version":1,"sha256Hash":"452f59168f3c5dacccc5fa161b5bf13d14e2cee2f6c7075f7f836cf4e695e4d7"}}`

	params := url.Values{}
	params.Set("operationName", "getProductResult")
	params.Set("variables", variables)
	params.Set("extensions", extensions)
	fullURL := "https://store.epicgames.com/graphql?" + params.Encode()

	fmt.Println(fullURL)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return RatingResponse{}, err
	}

	req.Header = http.Header{
		"accept":             {"*/*"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"sec-ch-ua":          {`"Google Chrome";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-fetch-dest":     {"empty"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"user-agent",
		},
	}

	resp, err := tlsClient.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return RatingResponse{}, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return RatingResponse{}, err
	}

	fmt.Println("response code:")
	fmt.Println(resp.StatusCode)
	fmt.Println("response body:")
	fmt.Println(string(responseData))
	var response RatingResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		fmt.Println("Error parsing response JSON:", string(responseData))
		return RatingResponse{}, err
	}
	return response, nil
}
