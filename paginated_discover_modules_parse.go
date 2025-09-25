package main

import (
	"encoding/json"
	"fmt"
)

// Main response structure
type PaginatedDiscoverModules struct {
	Typename *string           `json:"__typename,omitempty"`
	Data     []*DiscoverModule `json:"data,omitempty"`
	Paging   *PagingInfo       `json:"paging,omitempty"`
}

// Paging information
type PagingInfo struct {
	Count *int `json:"count,omitempty"`
	Start *int `json:"start,omitempty"`
	Total *int `json:"total,omitempty"`
}

// Main discover module structure
type DiscoverModule struct {
	Offers   []*Offer `json:"offers,omitempty"`
	Size     *string  `json:"size,omitempty"`
	TopicID  *string  `json:"topicId,omitempty"`
	Type     *string  `json:"type,omitempty"`
	Link     *Link    `json:"link,omitempty"`
	Title    *string  `json:"title,omitempty"`
}

// Link structure
type Link struct {
	LinkText *string `json:"linkText,omitempty"`
	Src      *string `json:"src,omitempty"`
}

// Offer structure
type Offer struct {
	Content   *Content `json:"content,omitempty"`
	OfferID   *string  `json:"offerId,omitempty"`
	SandboxID *string  `json:"sandboxId,omitempty"`
}

// Content structure
type Content struct {
	AgeRating     *AgeRating     `json:"ageRating,omitempty"`
	Attention     *Attention     `json:"attention,omitempty"`
	CatalogItemID *string        `json:"catalogItemId,omitempty"`
	Categories    []string       `json:"categories,omitempty"`
	Mapping       *Mapping       `json:"mapping,omitempty"`
	Media         *Media         `json:"media,omitempty"`
	Purchase      []*Purchase    `json:"purchase,omitempty"`
	Title         *string        `json:"title,omitempty"`
	SystemSpecs   *SystemSpecs   `json:"systemSpecs,omitempty"`
}

// Age rating structures
type AgeRating struct {
	AgeGate    *AgeGate        `json:"ageGate,omitempty"`
	AgeRating  *AgeRatingInfo  `json:"ageRating,omitempty"`
	RatingType *string         `json:"ratingType,omitempty"`
}

type AgeGate struct {
	Gate *string `json:"gate,omitempty"`
}

type AgeRatingInfo struct {
	AgeControl               *int      `json:"ageControl,omitempty"`
	ContentDescriptors       []string  `json:"contentDescriptors,omitempty"`
	InteractiveElements      []string  `json:"interactiveElements,omitempty"`
	RatingSystem             *string   `json:"ratingSystem,omitempty"`
	Title                    *string   `json:"title,omitempty"`
	RatingImage              *string   `json:"ratingImage,omitempty"`
	RectangularRatingImage   *string   `json:"rectangularRatingImage,omitempty"`
}

// Attention structure
type Attention struct {
	InAppPurchases *string `json:"inAppPurchases,omitempty"`
}

// Mapping structure
type Mapping struct {
	Slug *string `json:"slug,omitempty"`
}

// Media structures
type Media struct {
	AppIcon    *MediaItem `json:"appIcon,omitempty"`
	Card16x9   *MediaItem `json:"card16x9,omitempty"`
	Card3x4    *MediaItem `json:"card3x4,omitempty"`
	CoverImage *MediaItem `json:"coverImage,omitempty"`
}

type MediaItem struct {
	ImageSrc  *string `json:"imageSrc,omitempty"`
	AltText   *string `json:"altText,omitempty"`
	Height    *int    `json:"height,omitempty"`
	ImageType *string `json:"imageType,omitempty"`
	Width     *int    `json:"width,omitempty"`
}

// Purchase structures
type Purchase struct {
	OfferLifecycleState         *string          `json:"offerLifecycleState,omitempty"`
	Price                       *Price           `json:"price,omitempty"`
	PriceDisplay                *string          `json:"priceDisplay,omitempty"`
	PurchasePayload             *PurchasePayload `json:"purchasePayload,omitempty"`
	PurchaseStateEffectiveDate  *string          `json:"purchaseStateEffectiveDate,omitempty"`
	PurchaseType                *string          `json:"purchaseType,omitempty"`
	Discount                    *Discount        `json:"discount,omitempty"`
}

type Price struct {
	CurrencyCode *string  `json:"currencyCode,omitempty"`
	DecimalPrice *float64 `json:"decimalPrice,omitempty"`
}

type PurchasePayload struct {
	OfferID   *string `json:"offerId,omitempty"`
	SandboxID *string `json:"sandboxId,omitempty"`
}

type Discount struct {
	DiscountAmountDisplay *string `json:"discountAmountDisplay,omitempty"`
	OriginalPriceDisplay  *string `json:"originalPriceDisplay,omitempty"`
	DiscountEndDate       *string `json:"discountEndDate,omitempty"`
}

// System specs structures
type SystemSpecs struct {
	Platform             *string              `json:"platform,omitempty"`
	SystemRequirements   []*SystemRequirement `json:"systemRequirements,omitempty"`
	ApplicationVersion   *string              `json:"applicationVersion,omitempty"`
}

type SystemRequirement struct {
	RequirementType *string `json:"requirementType,omitempty"`
	Minimum         *string `json:"minimum,omitempty"`
	Title           *string `json:"title,omitempty"`
}

func ParsePaginatedDiscoverModules(body string) (*PaginatedDiscoverModules, error) {
	var result PaginatedDiscoverModules
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &result, nil
}
