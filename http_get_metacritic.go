package main

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	http "github.com/bogdanfinn/fhttp"
)

var reMetascore = regexp.MustCompile(`title="Metascore (\d+) out of 100"`)
var reMetacriticSlug = regexp.MustCompile(`href="/game/([^/"]+)/"`)

// toMetacriticSlug converts a game title to a Metacritic-style URL slug.
// e.g. "The Sims 4" → "the-sims-4", "Celeste" → "celeste"
func toMetacriticSlug(title string) string {
	title = strings.ToLower(title)
	var b strings.Builder
	prevHyphen := false
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevHyphen = false
		} else if !prevHyphen {
			b.WriteRune('-')
			prevHyphen = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// fetchMetacriticPage performs a tls-client GET to url and returns the body and status code.
func fetchMetacriticPage(rawURL string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header = http.Header{
		"accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"sec-ch-ua":          {`"Google Chrome";v="120", "Chromium";v="120", "Not=A?Brand";v="24"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"document"},
		"sec-fetch-mode":     {"navigate"},
		"sec-fetch-site":     {"none"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
		},
	}
	resp, err := tlsClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(data), resp.StatusCode, nil
}

// GetMetacriticScore searches Metacritic for gameTitle and returns the Metascore
// along with the Metacritic game slug.
//
// Strategy:
//  1. Fetch the search page to collect all /game/ slugs from search results.
//  2. Find the best-matching slug by comparing against the slugified game title:
//     prefer exact match, then prefix match, then fall back to the first result.
//  3. Fetch the game's own page (/game/{slug}/) to extract the score, since
//     older or less popular games may not show a Metascore badge in search listings.
func GetMetacriticScore(gameTitle string) (score int, metacriticSlug string, err error) {
	// Step 1: search for the game to find the best-matching slug.
	searchURL := "https://www.metacritic.com/search/" + url.PathEscape(gameTitle) + "/?category=13"
	fmt.Println("Search URL:", searchURL)

	body, statusCode, err := fetchMetacriticPage(searchURL)
	if err != nil {
		fmt.Println("Error fetching search page:", err)
		return 0, "", err
	}
	fmt.Println("Search response code:", statusCode)

	allSlugs := reMetacriticSlug.FindAllStringSubmatch(body, 20)
	titleSlug := toMetacriticSlug(gameTitle)
	fmt.Printf("Looking for slug matching: %q\n", titleSlug)

	// First pass: exact match.
	var prefixMatch string
	for _, m := range allSlugs {
		slug := m[1]
		if slug == titleSlug {
			metacriticSlug = slug
			break
		} else if prefixMatch == "" && strings.HasPrefix(slug, titleSlug+"-") {
			prefixMatch = slug
		}
	}
	// Second pass: prefix match if no exact found.
	if metacriticSlug == "" && prefixMatch != "" {
		metacriticSlug = prefixMatch
	}
	// Final fallback: first result.
	if metacriticSlug == "" && len(allSlugs) > 0 {
		metacriticSlug = allSlugs[0][1]
		fmt.Printf("No title match found; falling back to first result: %q\n", metacriticSlug)
	}

	if metacriticSlug == "" {
		fmt.Println("No game slugs found in search results")
		return 0, "", nil
	}

	// Step 2: fetch the game's own page for the score.
	gameURL := "https://www.metacritic.com/game/" + metacriticSlug + "/"
	fmt.Println("Game URL:", gameURL)

	gameBody, gameStatus, err := fetchMetacriticPage(gameURL)
	if err != nil {
		fmt.Println("Error fetching game page:", err)
		return 0, metacriticSlug, err
	}
	fmt.Println("Game page response code:", gameStatus)

	scoreMatches := reMetascore.FindStringSubmatch(gameBody)
	if len(scoreMatches) < 2 {
		fmt.Println("No Metascore found on game page")
		return 0, metacriticSlug, nil
	}

	score, err = strconv.Atoi(scoreMatches[1])
	if err != nil {
		fmt.Println("Error parsing score:", err)
		return 0, metacriticSlug, err
	}

	fmt.Printf("Metascore: %d (slug: %s)\n", score, metacriticSlug)
	return score, metacriticSlug, nil
}
