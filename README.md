# Building
```
go build main.go .\cli_hander_free.go .\cli_hander_rating.go .\cli_handler_search.go .\graphql_get_free_games.go .\graphql_get_rating.go .\graphql_search_for_game_info.go
```

# Running


## Search
Searching based on input json
```
.\main search --inputFile .\epic_free_games.json --outputFile out.json > search.log 2>&1
```

Single search:
```
.\main search --gameTitle "Celeste"
```

## Ratings
Rating based on input json
```
.\main rate --inputFile .\epic_free_games.json --outputFile out.json > ratings.log 2>&1
```

Single rating:
```
.\main rate --searchKey b671fbc7be424e888c9346a9a6d3d9db 
```

## Free Games


Append to file:
```
.\main free --inputFile epic_free_games.json --outputFile out.json
```

Print to console:
```
.\main free 
```

# High Level Solution

1. Start with a list of game title for free games.
2. Use epic's graphql api to search for the url and sandboxId by using the game title (`searchStoreQuery`)
    1. use mapping's slug if there
    2. otherwise use productSlug with /home removed from the end
    3. otherwise use urlSlug
    4. manually resolve and weirdness
3. Use epic's graphql api to get the rating using the epic game's namespace as the sandboxId (`getProductResult`)

# Research

## chromedp failure
```shell
go get -u github.com/chromedp/chromedp
```
chromedp was unsuccessful since cloudflare gives us a challenge which we cannot answer.

## stackoverflow to the rescue
instead I found this stack overflow
https://stackoverflow.com/questions/75594997/is-it-possible-to-get-data-about-specific-game-from-epic-games-store

I was able to run the query without issue in powershell:
```powershell
$response = Invoke-WebRequest -Uri 'https://graphql.epicgames.com/graphql?query=query%20searchStoreQuery(%20$allowCountries:%20String%20$category:%20String%20$namespace:%20String%20$itemNs:%20String%20$sortBy:%20String%20$sortDir:%20String%20$start:%20Int%20$tag:%20String%20$releaseDate:%20String%20$withPrice:%20Boolean%20=%20true%20)%20{%20Catalog%20{%20searchStore(%20allowCountries:%20$allowCountries%20category:%20$category%20count:%201%20country:%20%22US%22%20keywords:%20%22celeste%22%20namespace:%20$namespace%20itemNs:%20$itemNs%20sortBy:%20$sortBy%20sortDir:%20$sortDir%20releaseDate:%20$releaseDate%20start:%20$start%20tag:%20$tag%20)%20{%20elements%20{%20title%20description%20keyImages%20{%20type%20url%20}%20seller%20{%20name%20}%20categories%20{%20path%20}%20price(country:%20%22US%22)%20@include(if:%20$withPrice)%20{%20totalPrice%20{%20fmtPrice(locale:%20%22en-US%22)%20{%20discountPrice%20}%20}%20}%20}%20}%20}%20}'
Write-Output $response.StatusCode
200
Write-Output $response.Content
{"data":{"Catalog":{"searchStore":{"elements":[{"title":"Celeste","description":"Help Madeline survive her inner demons on her journey to the top of Celeste Mountain, in this super-tight platformer from the creators of TowerFall. Brave hundreds of hand-crafted challenges, uncover devious secrets, and piece together the mystery of the mountain.","keyImages":[{"type":"OfferImageWide","url":"https://cdn1.epicgames.com/b671fbc7be424e888c9346a9a6d3d9db/offer/Celeste - landscape offer image-2560x1440-0b9b94fd493d817704ecfdf4c704989a.jpg"},{"type":"OfferImageTall","url":"https://cdn1.epicgames.com/offer/b671fbc7be424e888c9346a9a6d3d9db/CodeRedemption_Celeste-340x440-873dacb76be8c59d659757b4b5284a611_1200x1600-9e39079f1ce19738e65d349a8bc98650"},{"type":"Thumbnail","url":"https://cdn1.epicgames.com/offer/b671fbc7be424e888c9346a9a6d3d9db/CodeRedemption_Celeste-340x440-873dacb76be8c59d659757b4b5284a611_1200x1600-9e39079f1ce19738e65d349a8bc98650"}],"seller":{"name":"Maddy Makes Games"},"categories":[{"path":"games"},{"path":"games/edition/base"},{"path":"games/edition"},{"path":"applications"}],"price":{"totalPrice":{"fmtPrice":{"discountPrice":"$19.99"}}}}]}}},"extensions":{}}
```

next step is to explore their graphql API and look for the rating some where.




## Searching:
```graphql
query searchStoreQuery($allowCountries: String, $category: String, $count: Int, $country: String!, $keywords: String, $locale: String, $namespace: String, $itemNs: String, $sortBy: String, $sortDir: String, $start: Int, $tag: String, $releaseDate: String, $withPrice: Boolean = false, $withPromotions: Boolean = false, $priceRange: String, $freeGame: Boolean, $onSale: Boolean, $effectiveDate: String) {
  Catalog {
    searchStore(
      allowCountries: $allowCountries
      category: $category
      count: $count
      country: $country
      keywords: $keywords
      locale: $locale
      namespace: $namespace
      itemNs: $itemNs
      sortBy: $sortBy
      sortDir: $sortDir
      releaseDate: $releaseDate
      start: $start
      tag: $tag
      priceRange: $priceRange
      freeGame: $freeGame
      onSale: $onSale
      effectiveDate: $effectiveDate
    ) {
      elements {
        title
        id
        namespace
        description
        effectiveDate
        keyImages {
          type
          url
        }
        currentPrice
        seller {
          id
          name
        }
        productSlug
        urlSlug
        url
        tags {
          id
        }
        items {
          id
          namespace
        }
        customAttributes {
          key
          value
        }
        categories {
          path
        }
        catalogNs {
          mappings(pageType: "productHome") {
            pageSlug
            pageType
          }
        }
        offerMappings {
          pageSlug
          pageType
        }
        price(country: $country) @include(if: $withPrice) {
          totalPrice {
            discountPrice
            originalPrice
            voucherDiscount
            discount
            currencyCode
            currencyInfo {
              decimals
            }
            fmtPrice(locale: $locale) {
              originalPrice
              discountPrice
              intermediatePrice
            }
          }
          lineOffers {
            appliedRules {
              id
              endDate
              discountSetting {
                discountType
              }
            }
          }
        }
        promotions(category: $category) @include(if: $withPromotions) {
          promotionalOffers {
            promotionalOffers {
              startDate
              endDate
              discountSetting {
                discountType
                discountPercentage
              }
            }
          }
          upcomingPromotionalOffers {
            promotionalOffers {
              startDate
              endDate
              discountSetting {
                discountType
                discountPercentage
              }
            }
          }
        }
      }
      paging {
        count
        total
      }
    }
  }
}
```
https://github.com/Tectors/EpicGraphQL/blob/main/docs/graphql/catalog/searchStoreQuery.graphql

### Search for useful info like id, url:
```
https://graphql.epicgames.com/graphql?query=query  searchStoreQuery($allowCountries: String, $category: String, $locale: String, $namespace: String, $itemNs: String, $sortBy: String, $sortDir: String, $start: Int, $tag: String, $releaseDate: String, $priceRange: String, $freeGame: Boolean, $onSale: Boolean, $effectiveDate: String) {
  Catalog {
    searchStore(
      allowCountries: $allowCountries
      category: $category
      count: 1
      country: "US"
      keywords: "Wonder Boy"
      locale: $locale
      namespace: $namespace
      itemNs: $itemNs
      sortBy: $sortBy
      sortDir: $sortDir
      releaseDate: $releaseDate
      start: $start
      tag: $tag
      priceRange: $priceRange
      freeGame: $freeGame
      onSale: $onSale
      effectiveDate: $effectiveDate
    ) {
      elements {
        title
        id
        namespace
        description
        effectiveDate
        productSlug
        urlSlug
        url
                tags {
          id
        }
        items {
          id
          namespace
        }
        customAttributes {
          key
          value
        }
        categories {
          path
        }
        catalogNs {
          mappings(pageType: "productHome") {
            pageSlug
            pageType
          }
        }
        offerMappings {
          pageSlug
          pageType
        }
      }
    }
  }
}
```

For `Celeste` it's in 
```
data.Catalog.searchStore.elements[0].urlSlug
data.Catalog.searchStore.elements[0].catalogNs.mappings[0].pageSlug
```

<details>
<summary>Click to see response</summary>
```
{
  "data": {
    "Catalog": {
      "searchStore": {
        "elements": [
          {
            "title": "Celeste",
            "id": "544c18ccadf8466691d8f2a335721b14",
            "namespace": "b671fbc7be424e888c9346a9a6d3d9db",
            "description": "Help Madeline survive her inner demons on her journey to the top of Celeste Mountain, in this super-tight platformer from the creators of TowerFall. Brave hundreds of hand-crafted challenges, uncover devious secrets, and piece together the mystery of the mountain.",
            "effectiveDate": "2019-08-29T15:00:00.000Z",
            "productSlug": "celeste/home",
            "urlSlug": "celeste",
            "url": null,
            "tags": [
              {
                "id": "1216"
              },
              {
                "id": "1188"
              },
              {
                "id": "21894"
              },
              {
                "id": "21129"
              },
              {
                "id": "9547"
              },
              {
                "id": "9549"
              },
              {
                "id": "21134"
              },
              {
                "id": "1263"
              },
              {
                "id": "21138"
              },
              {
                "id": "21139"
              },
              {
                "id": "21140"
              },
              {
                "id": "21109"
              },
              {
                "id": "21141"
              },
              {
                "id": "1370"
              },
              {
                "id": "21149"
              },
              {
                "id": "10719"
              },
              {
                "id": "1151"
              }
            ],
            "items": [
              {
                "id": "1e6af8b98bb644d7ac51fc810bc36d85",
                "namespace": "b671fbc7be424e888c9346a9a6d3d9db"
              }
            ],
            "customAttributes": [
              {
                "key": "com.epicgames.app.blacklist",
                "value": "{}"
              },
              {
                "key": "com.epicgames.app.productSlug",
                "value": "celeste/home"
              }
            ],
            "categories": [
              {
                "path": "games"
              },
              {
                "path": "games/edition/base"
              },
              {
                "path": "games/edition"
              },
              {
                "path": "applications"
              }
            ],
            "catalogNs": {
              "mappings": [
                {
                  "pageSlug": "celeste",
                  "pageType": "productHome"
                }
              ]
            },
            "offerMappings": []
          }
        ]
      }
    }
  },
  "extensions": {}
}
```
</details>


For `Wonder Boy` it's in
```
data.Catalog.searchStore.elements[0].catalogNs.mappings[0].pageSlug
```

<details>
<summary>Click to see response</summary>
```
{
  "data": {
    "Catalog": {
      "searchStore": {
        "elements": [
          {
            "title": "Wonder Boy The Dragons Trap",
            "id": "bd18a76d848946b0ac5f5c692c8b4757",
            "namespace": "e1e92771f6774072bb1b8d0a0a6197f7",
            "description": "Boasting beautiful, hand-drawn animations and a re-orchestrated soundtrack, the cult classic returns with a unique blend of exploration, action, and adventure!",
            "effectiveDate": "2022-07-14T15:00:00.000Z",
            "productSlug": null,
            "urlSlug": "486c4f8c133c4930a295ecffe84a80d4",
            "url": null,
            "tags": [
              {
                "id": "21122"
              },
              {
                "id": "1188"
              },
              {
                "id": "21127"
              },
              {
                "id": "9547"
              },
              {
                "id": "1263"
              },
              {
                "id": "21138"
              },
              {
                "id": "21139"
              },
              {
                "id": "21140"
              },
              {
                "id": "21141"
              },
              {
                "id": "1336"
              },
              {
                "id": "1370"
              },
              {
                "id": "21149"
              },
              {
                "id": "1151"
              },
              {
                "id": "21119"
              }
            ],
            "items": [
              {
                "id": "da28acf70b264b1eb8cc3e2a74bdc530",
                "namespace": "e1e92771f6774072bb1b8d0a0a6197f7"
              }
            ],
            "customAttributes": [
              {
                "key": "autoGeneratedPrice",
                "value": "false"
              },
              {
                "key": "isManuallySetPCReleaseDate",
                "value": "true"
              }
            ],
            "categories": [
              {
                "path": "games/edition/base"
              },
              {
                "path": "games/edition"
              },
              {
                "path": "games"
              }
            ],
            "catalogNs": {
              "mappings": [
                {
                  "pageSlug": "wonder-boy-the-dragons-trap-26381d",
                  "pageType": "productHome"
                }
              ]
            },
            "offerMappings": [
              {
                "pageSlug": "wonder-boy-the-dragons-trap-26381d",
                "pageType": "productHome"
              }
            ]
          }
        ]
      }
    }
  },
  "extensions": {}
}
```
</details>

As a result, I'm going to use `data.Catalog.searchStore.elements[0].catalogNs.mappings[0].pageSlug` for getting the url.

After writing some code and testing,

sometimes we have `data.Catalog.searchStore.elements[0].catalogNs.mappings[0].pageSlug`

sometimes we have `data.Catalog.searchStore.elements[0].urlSlug`

sometimes we have both

## product review
```graphql
query productReviewsQuery($sku: String!) {
    OpenCritic {
        productReviews(sku: $sku) {
            id
            name
            openCriticScore
            reviewCount
            percentRecommended
            openCriticUrl
            award
            topReviews {
                publishedDate
                externalUrl
                snippet
                language
                score
                author
                ScoreFormat {
                    id
                    description
                }
                OutletId
                outletName
                displayScore
            }
        }
    }
}
```
https://github.com/SD4RK/epicstore_api/blob/master/epicstore_api/queries.py



```
https://graphql.epicgames.com/graphql?query=query productReviewsQuery($sku: String! = "6000206130537") {
    OpenCritic {
        productReviews(sku: $sku) {
            id
            name
            openCriticScore
            reviewCount
            percentRecommended
            openCriticUrl
            award
            topReviews {
                publishedDate
                externalUrl
                snippet
                language
                score
                author
                ScoreFormat {
                    id
                    description
                }
                OutletId
                outletName
                displayScore
            }
        }
    }
}
```
tried ID from Celesete but didn't get anything back.


By searching for 4.9 (rating of Celeste, I found in HTTP GET of game page there is:
```
                    "state": {
                        "data": {
                            "RatingsPolls": {
                                "getProductResult": {
                                    "averageRating": 4.92,
                                    "pollResult": [{
                                        "id": 65,
                                        "tagId": 21109,
                                        "pollDefinitionId": 1,
                                        "localizations": {
                                            "text": "Yes",
                                            "emoji": "https:\u002F\u002Fcdn2.epicgames.com\u002Fstatic\u002Ffonts\u002Fjoypixel\u002F2705.svg",
                                            "resultEmoji": "https:\u002F\u002Fcdn2.epicgames.com\u002Fstatic\u002Ffonts\u002Fjoypixel\u002F1f996.svg",
                                            "resultTitle": "Great Boss Battles",
                                            "resultText": "This game has"
                                        },
                                        "total": 1655
                                    }, {
```
which is exactly what I want and looks like graphql. Going to play around with RatingsPolls.


```
      "RatingsPolls": {
          "getProductResult": {
              "averageRating": 4.92,
              "pollResult": [{
```

id didn't work, but namespace did:
```
https://graphql.epicgames.com/graphql?query=query getProductResult($sandboxId: String = "b671fbc7be424e888c9346a9a6d3d9db", $locale: String = "US") {
  RatingsPolls {
    getProductResult(sandboxId: $sandboxId, locale: $locale) {
      averageRating
    }
  }
}
```
gave back
```
{"data":{"RatingsPolls":{"getProductResult":{"averageRating":4.92}}},"extensions":{}}
```

lets try a game with a different rating just to double check.

```
https://graphql.epicgames.com/graphql?query=query getProductResult($sandboxId: String = "e1e92771f6774072bb1b8d0a0a6197f7", $locale: String = "US") {
  RatingsPolls {
    getProductResult(sandboxId: $sandboxId, locale: $locale) {
      averageRating
    }
  }
}
```
gave back
```
{"data":{"RatingsPolls":{"getProductResult":{"averageRating":4.44}}},"extensions":{}}
```
as expected by browsing to the Wonder Boy product page

We have all the information we need to write a program to populate everything.