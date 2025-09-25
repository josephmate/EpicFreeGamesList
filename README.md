todo try:
1. worked https://github.com/bogdanfinn/tls-client
2. did not try https://github.com/enetx/surf
3. did not work github.com/BridgeSenseDev/go-curl-impersonate

Every Thursday 19:00 UTC, gets the latest free game from the Epic Store and adds it to /epic_free_games.json . That file is used to render the free games table on https://josephmate.github.io/EpicFreeGamesList/ .

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


# Making docker image

```
docker build -t josephmate/epic-games-free-list-updater .
```

# Sending to Dockerhub

```
docker push josephmate/epic-games-free-list-updater
```

# Running from docker image


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


## getting the real snapshotId

There are a lot of 0s. Initially I attributed that to those pages not having ratings. However, [MCHopie
 noticed that a majority of the pages with 0 acutally had ratings.](https://www.reddit.com/r/EpicGamesPC/comments/1anepsx/epic_free_games_list_with_sorting_by_epic_ratings/ks3oedi/?context=8&depth=9).
 
 With deeper investigation in a particular example Payday 2, the snapshotId used to obtain the rating was wrong. In this case I had `d5241c76f178492ea1540fce45616757` which is not found in the rating API. When you visit the Payday 2 page and look at the javascript source code they have this embedded in it:

 ```
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
                                        "offerId": "384f75fdd6b34f63a2daace1a3c5dab0",
                                        "offer": {
                                            "id": "384f75fdd6b34f63a2daace1a3c5dab0",
                                            "namespace": "3b661fd6a9724ac7b6ac6d10d0572511",
                                            "effectiveDate": "2023-06-08T15:00:00.000Z",
                                            "expiryDate": null
                                        },
                                        "prePurchaseOfferId": null,
                                        "prePurchaseOffer": null,
                                        "pageId": null
                                    }
                                }
                            }
                        },
 ```

 Notice two things. First the snapshotId is different: `3b661fd6a9724ac7b6ac6d10d0572511`. Second, the javascript seems to be pre-populated with the GraphQL responses. Unfortunately, I do not know what query maps to result `StorePageMapping`. Fortunately, [woctezuma had a page that dumped all the GraphQL Queries](https://gist.github.com/woctezuma/8ca464a276b15d7dfad475fd6b6cbee9):

 ```
query getMappingByPageSlug($pageSlug: String!, $sandboxId: String, $locale: String!) {
  StorePageMapping {
    mapping(pageSlug: $pageSlug, sandboxId: $sandboxId) {
      pageSlug
      pageType
      sandboxId
      productId
      createdDate
      updatedDate
      mappings {
        cmsSlug
        offerId
        offer(locale: $locale) {
          id
          namespace
          effectiveDate
          expiryDate
        }
        prePurchaseOfferId
        prePurchaseOffer(locale: $locale) {
          id
          namespace
          effectiveDate
          expiryDate
        }
        pageId
      }
    }
  }
}
```



```
https://graphql.epicgames.com/graphql?query=query getProductResult($sandboxId: String = "b671fbc7be424e888c9346a9a6d3d9db", $locale: String = "US") {
  RatingsPolls {
    getProductResult(sandboxId: $sandboxId, locale: $locale) {
      averageRating
    }
  }
}
```


Here's my old information for Payday 2:
```
  {
    "epicId": "de434b7be57940d98ede93b50cdacfc2",
    "epicRating": 0,
    "epicStoreLink": "https://store.epicgames.com/en-US/p/payday-2-c66369",
    "freeDate": "2023-06-08",
    "gameTitle": "Payday 2",
    "mappingSlug": "",
    "productSlug": "payday-2-c66369",
    "sandboxId": "d5241c76f178492ea1540fce45616757",
    "urlSlug": "mystery-game-7"
  },
```

```
query getMappingByPageSlug($pageSlug: String! = "payday-2-c66369", $sandboxId: String) {
  StorePageMapping {
    mapping(pageSlug: $pageSlug, sandboxId: $sandboxId) {
      pageSlug
      pageType
      sandboxId
      productId
      createdDate
      updatedDate
      mappings {
        cmsSlug
        pageId
      }
    }
  }
}
```


```powershell
$response = Invoke-WebRequest -Uri 'https://graphql.epicgames.com/graphql?query=query%20getMappingByPageSlug%28%24pageSlug%3A%20String%21%20%3D%20%22payday-2-c66369%22%2C%20%24sandboxId%3A%20String%29%20%7B%0A%20%20StorePageMapping%20%7B%0A%20%20%20%20mapping%28pageSlug%3A%20%24pageSlug%2C%20sandboxId%3A%20%24sandboxId%29%20%7B%0A%20%20%20%20%20%20pageSlug%0A%20%20%20%20%20%20pageType%0A%20%20%20%20%20%20sandboxId%0A%20%20%20%20%20%20productId%0A%20%20%20%20%20%20createdDate%0A%20%20%20%20%20%20updatedDate%0A%20%20%20%20%20%20mappings%20%7B%0A%20%20%20%20%20%20%20%20cmsSlug%0A%20%20%20%20%20%20%20%20pageId%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D'
Write-Output $response.StatusCode
200
Write-Output $response.Content
{"data":{"StorePageMapping":{"mapping":{"pageSlug":"payday-2-c66369","pageType":"productHome","sandboxId":"3b661fd6a9724ac7b6ac6d10d0572511","productId":"14eb3477a6084940b49de5aa73c60f98","createdDate":"2023-06-07T08:05:53.761Z","updatedDate":"2023-11-27T20:21:20.490Z","mappings":{"cmsSlug":null,"pageId":null}}}},"extensions":{}}
```

Notice that the sandboxId returned matches the one I need to get the ratings: `3b661fd6a9724ac7b6ac6d10d0572511`.

Trying again with another example `Death Stranding`, with productSlug `death-stranding%2Fhome`:
```
$response = Invoke-WebRequest -Uri 'https://graphql.epicgames.com/graphql?query=query%20getMappingByPageSlug%28%24pageSlug%3A%20String%21%20%3D%20%22death-stranding%2Fhome%22%2C%20%24sandboxId%3A%20String%29%20%7B%0A%20%20StorePageMapping%20%7B%0A%20%20%20%20mapping%28pageSlug%3A%20%24pageSlug%2C%20sandboxId%3A%20%24sandboxId%29%20%7B%0A%20%20%20%20%20%20pageSlug%0A%20%20%20%20%20%20pageType%0A%20%20%20%20%20%20sandboxId%0A%20%20%20%20%20%20productId%0A%20%20%20%20%20%20createdDate%0A%20%20%20%20%20%20updatedDate%0A%20%20%20%20%20%20mappings%20%7B%0A%20%20%20%20%20%20%20%20cmsSlug%0A%20%20%20%20%20%20%20%20pageId%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D'
Write-Output $response.StatusCode
200
Write-Output $response.Content
{"data":{"StorePageMapping":{"mapping":null}},"extensions":{}}
```
didn't work. but urlSlug did `death-stranding` :
```
$response = Invoke-WebRequest -Uri 'https://graphql.epicgames.com/graphql?query=query%20getMappingByPageSlug%28%24pageSlug%3A%20String%21%20%3D%20%22death-stranding%22%2C%20%24sandboxId%3A%20String%29%20%7B%0A%20%20StorePageMapping%20%7B%0A%20%20%20%20mapping%28pageSlug%3A%20%24pageSlug%2C%20sandboxId%3A%20%24sandboxId%29%20%7B%0A%20%20%20%20%20%20pageSlug%0A%20%20%20%20%20%20pageType%0A%20%20%20%20%20%20sandboxId%0A%20%20%20%20%20%20productId%0A%20%20%20%20%20%20createdDate%0A%20%20%20%20%20%20updatedDate%0A%20%20%20%20%20%20mappings%20%7B%0A%20%20%20%20%20%20%20%20cmsSlug%0A%20%20%20%20%20%20%20%20pageId%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D'
Write-Output $response.Content
{"data":{"StorePageMapping":{"mapping":{"pageSlug":"death-stranding","pageType":"productHome","sandboxId":"f4a904fcef2447439c35c4e6457f3027","productId":"da519d41698b4854815db7371210c3a1","createdDate":"2021-05-05T16:55:53.681Z","updatedDate":"2023-05-31T16:12:45.187Z","mappings":{"cmsSlug":"death-stranding/home","pageId":null}}}},"extensions":{}}
```

`f4a904fcef2447439c35c4e6457f3027` matches the one on the Death Standing page!
