package main

import (
	"errors"
	"fmt"
	"os"
)

func getGames(url string) (*PaginatedDiscoverModules, error) {
	result, err := HttpGetPaginatedDiscoverModules(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching discover modules: %v\n", err)
		return nil, err
	}

	// has non-empty string
	if len(result) > 0 {
		parsedDiscoverModules, err := ParsePaginatedDiscoverModules(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing discover modules: %v\n", err)
			return nil, err
		}

		return parsedDiscoverModules, nil
	} else {
		return nil, errors.New("empty response body")
	}
}

func notFree(purchases []*Purchase) (bool) {
	for _, purchase := range purchases {
		if (purchase.Discount != nil && purchase.Discount.DiscountAmountDisplay != nil && *purchase.Discount.DiscountAmountDisplay == "-100%") {
			return false
		}
	}

	return true
}

func getFreeGames(url string, platform string) ([]FreeGameEntry, error) {
	games, err := getGames(url)
	if (err != nil) {
			fmt.Fprintf(os.Stderr, "Failed to get games: %v\n", err)
			return nil, err
	}

	if (games == nil) {
			return nil, errors.New("nil PaginatedDiscoverModules")
	}

	
	var freeGames = []FreeGameEntry{}
	for _, game := range games.Data {
		if (game.Offers == nil) {
			continue
		}
		for _, offer := range game.Offers {
			if offer.Content == nil {
				continue
			}

			if offer.Content.Mapping == nil || offer.Content.Mapping.Slug == nil {
				continue
			}

			if offer.Content.Title == nil {
				continue
			}

			if notFree(offer.Content.Purchase) {
				continue
			}

			slugId := offer.Content.Mapping.Slug
			freeGame := FreeGameEntry{}
			freeGame.EpicStoreLink = "https://store.epicgames.com/en-US/p/" + *slugId
			freeGame.MappingSlug = *slugId
			freeGame.Platform = platform
			freeGame.GameTitle = *offer.Content.Title

			freeGames = append(freeGames, freeGame)
		}
	}

	return freeGames, nil
}

func FreeMobileGames() (*FreeGames, error) {

  var allFreeGames = []FreeGameEntry{}

	androidGames, err := getFreeGames(PaginatedDiscoverModulesAndroidURL, "android")
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Error retrieving free android games: %v\n", err)
		return nil, err
	}
	iosGames, err := getFreeGames(PaginatedDiscoverModulesIosURL, "ios")
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Error retrieving free ios games: %v\n", err)
		return nil, err
	}

	allFreeGames = append(allFreeGames, androidGames...)
	allFreeGames = append(allFreeGames, iosGames...)

	result := FreeGames{}
	result.ThisWeek = allFreeGames

	return &result, nil
}
