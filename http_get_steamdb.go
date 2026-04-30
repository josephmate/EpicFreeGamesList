package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
)

type steamSearchResponse struct {
	Total int `json:"total"`
	Items []struct {
		Type string `json:"type"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"items"`
}

type steamReviewsResponse struct {
	Success      int `json:"success"`
	QuerySummary struct {
		TotalPositive int `json:"total_positive"`
		TotalReviews  int `json:"total_reviews"`
	} `json:"query_summary"`
}

func fetchSteamJSON(rawURL string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

// steamDBRatingFormula computes the SteamDB rating from review counts.
// Formula: Rating = ReviewScore - (ReviewScore - 0.5) * 2^(-log10(TotalReviews + 1))
// Returns a percentage (0–100), rounded to 2 decimal places.
func steamDBRatingFormula(totalPositive, totalReviews int) float64 {
	if totalReviews == 0 {
		return 0
	}
	reviewScore := float64(totalPositive) / float64(totalReviews)
	raw := reviewScore - (reviewScore-0.5)*math.Pow(2, -math.Log10(float64(totalReviews)+1))
	return math.Round(raw*10000) / 100
}

// GetSteamDBRating searches Steam for gameTitle, picks the best-matching app,
// fetches its review summary, and returns the SteamDB-style rating (0–100 percentage)
// along with the Steam App ID. Returns (0, 0, nil) if no match is found.
func GetSteamDBRating(gameTitle string) (rating float64, steamAppID int, err error) {
	searchURL := "https://store.steampowered.com/api/storesearch?term=" + url.QueryEscape(gameTitle) + "&l=english&cc=US"
	fmt.Println("Steam search URL:", searchURL)

	var searchResult steamSearchResponse
	if err = fetchSteamJSON(searchURL, &searchResult); err != nil {
		return 0, 0, fmt.Errorf("steam search failed: %w", err)
	}
	fmt.Printf("Steam search: %d results\n", searchResult.Total)

	// Find best app match by bigram similarity (reuses bigramSim from http_get_metacritic.go).
	// Only consider type=="app" to exclude soundtracks, videos, etc.
	normalizedTitle := strings.ToLower(strings.TrimSpace(gameTitle))
	primaryT := strings.ToLower(primaryTitle(gameTitle))
	bestID := 0
	bestScore := -1.0
	for _, item := range searchResult.Items {
		if item.Type != "app" {
			continue
		}
		candidate := strings.ToLower(item.Name)
		// Exact match (full title or primary title) wins immediately.
		if candidate == normalizedTitle || candidate == primaryT {
			bestID = item.ID
			bestScore = 1.0
			fmt.Printf("Steam exact match: %q (appID=%d)\n", item.Name, item.ID)
			break
		}
		score := bigramSim(normalizedTitle, candidate)
		if score > bestScore {
			bestScore = score
			bestID = item.ID
		}
	}

	if bestID == 0 {
		fmt.Println("No Steam app found for:", gameTitle)
		return 0, 0, nil
	}
	if bestScore < 1.0 {
		fmt.Printf("Steam fuzzy match: appID=%d (score=%.3f)\n", bestID, bestScore)
	}

	reviewURL := fmt.Sprintf("https://store.steampowered.com/appreviews/%d?json=1&language=all&purchase_type=all&num_per_page=0", bestID)
	fmt.Println("Steam reviews URL:", reviewURL)

	var reviewResult steamReviewsResponse
	if err = fetchSteamJSON(reviewURL, &reviewResult); err != nil {
		return 0, bestID, fmt.Errorf("steam reviews fetch failed: %w", err)
	}
	if reviewResult.Success != 1 {
		fmt.Println("Steam reviews API returned failure")
		return 0, bestID, nil
	}

	totalPositive := reviewResult.QuerySummary.TotalPositive
	totalReviews := reviewResult.QuerySummary.TotalReviews
	fmt.Printf("Reviews: %d positive / %d total\n", totalPositive, totalReviews)

	if totalReviews == 0 {
		fmt.Println("No reviews found")
		return 0, bestID, nil
	}

	rating = steamDBRatingFormula(totalPositive, totalReviews)
	fmt.Printf("SteamDB Rating: %.2f%%\n", rating)
	return rating, bestID, nil
}
