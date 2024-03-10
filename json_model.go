package main

/*
  {
    "epicId": "014f265f264f46e6b5d59c738cf24ee4",
    "epicRating": 4.62,
    "epicStoreLink": "https://www.epicgames.com/store/p/the-sims-4",
    "freeDate": "2023-05-11",
    "gameTitle": "The Sims™ 4 The Daring Lifestyle Bundle",
    "mappingSlug": "the-sims-4",
    "productSlug": "the-sims-4",
    "sandboxId": "2a14cf8a83b149919a2399504e5686a6",
    "urlSlug": "the-sims-4"
  },
*/
type GameEntryComplete struct {
	EpicId        string `json:"epicId"`
	EpicRating    float64 `json:"epicRating"`
	EpicStoreLink string `json:"epicStoreLink"`
	FreeDate      string `json:"freeDate"`
	GameTitle     string `json:"gameTitle"`
	MappingSlug   string `json:"mappingSlug"`
	ProductSlug   string `json:"productSlug"`
	SandboxId     string `json:"sandboxId"`
	UrlSlug       string `json:"urlSlug"`
}
