package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FreeGameEntry struct {
	EpicId        string `json:"epicId"`
	EpicStoreLink string `json:"epicStoreLink"`
	FreeDate      string `json:"freeDate"`
	GameTitle     string `json:"gameTitle"`
	MappingSlug   string `json:"mappingSlug"`
	ProductSlug   string `json:"productSlug"`
	SandboxId     string `json:"sandboxId"`
	UrlSlug       string `json:"urlSlug"`
}

type FreeGames struct {
	ThisWeek []FreeGameEntry
	NextWeek []FreeGameEntry
}

type DiscountSetting struct {
	DiscountType       *string `json:"discountType"`
	DiscountPercentage *int    `json:"discountPercentage"`
}

type PromotionalOffer struct {
	PromotionalOffers []struct {
		DiscountSetting DiscountSetting `json:"discountSetting"`
		StartDate       string          `json:"startDate"`
		EndDate         string          `json:"endDate"`
	} `json:"promotionalOffers"`
}

type GameEntry struct {
	CatalogNs struct {
		Mappings []struct {
			PageSlug string `json:"pageSlug"`
		} `json:"mappings"`
	} `json:"catalogNs"`
	Categories []struct {
		Path string `json:"path"`
	} `json:"categories"`
	Namespace   string `json:"namespace"`
	Id          string `json:"id"`
	ProductSlug string `json:"productSlug"`
	Promotions  struct {
		PromotionalOffers         []PromotionalOffer `json:"promotionalOffers"`
		UpcomingPromotionalOffers []PromotionalOffer `json:"upcomingPromotionalOffers"`
	} `json:"promotions"`
	Seller struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"seller"`
	UrlSlug string `json:"urlSlug"`
}

type FreeGameResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []GameEntry `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

const epicDevAccountId = "o-ufmrk5furrrxgsp5tdngefzt5rxdcn"
const epicDevAccountName = "Epic Dev Test Account"

func isVaultedGame(game GameEntry) bool {
	seller := game.Seller

	if seller.Id == epicDevAccountId {
		return true
	}
	if seller.Name == epicDevAccountName {
		return true
	}

	categories := game.Categories

	if len(categories) > 0 {
		for _, category := range categories {
			if category.Path == "freegames/vaulted" {
				return true
			}
		}
	}

	return false
}

// sometimes a game has multiple promotions going on. see sample_free_game.json
// for an example
func anyFree(promotionalOffers []PromotionalOffer) *string {
	for _, promotionalOfferOuter := range promotionalOffers {
		for _, promotionalOffersInner := range promotionalOfferOuter.PromotionalOffers {
			if promotionalOffersInner.DiscountSetting.DiscountType != nil &&
				*promotionalOffersInner.DiscountSetting.DiscountType == "PERCENTAGE" &&
				promotionalOffersInner.DiscountSetting.DiscountPercentage != nil &&
				*promotionalOffersInner.DiscountSetting.DiscountPercentage == 0 {
				return &strings.Split(promotionalOffersInner.StartDate, "T")[0]
			}
		}
	}

	return nil

}

func GetFreeGames() (FreeGames, error) {

	// Perform the HTTP GET request to the GraphQL endpoint
	resp, err := http.Get("https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=CA&allowCountries=CA")
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return FreeGames{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return FreeGames{}, err
	}
	//fmt.Println("response:\n", string(responseData))

	var response FreeGameResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		fmt.Println("Error parsing response JSON:", string(responseData))
		return FreeGames{}, err
	}

	result := FreeGames{}
	for _, element := range response.Data.Catalog.SearchStore.Elements {
		if !isVaultedGame(element) {
			modifiedEntry := FreeGameEntry{}
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

			thisWeekFree := anyFree(element.Promotions.PromotionalOffers)
			nextWeekFree := anyFree(element.Promotions.UpcomingPromotionalOffers)
			if thisWeekFree != nil {
				modifiedEntry.FreeDate = *thisWeekFree
				result.ThisWeek = append(result.ThisWeek, modifiedEntry)
			} else if nextWeekFree != nil {
				modifiedEntry.FreeDate = *nextWeekFree
				result.NextWeek = append(result.NextWeek, modifiedEntry)
			}
		}
	}

	return result, nil
}
