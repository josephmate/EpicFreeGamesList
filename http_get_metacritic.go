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
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var reMetascore = regexp.MustCompile(`data-testid="global-score-value">(\d+)<`)
var reMetacriticSlug = regexp.MustCompile(`href="/game/([^/"]+)/"`)

// toASCII converts accented/decorated characters to their ASCII base letter
// using NFD decomposition + stripping combining marks (e.g. û→u, é→e, ü→u).
func toASCII(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// toMetacriticSlug converts a game title to a Metacritic-style URL slug.
//   - ASCII-normalizes accented chars (û→u)
//   - Strips apostrophes/curly-apostrophes so "Wake's"→"wakes" not "wake-s"
//   - Strips dots so "Q.U.B.E."→"qube" not "q-u-b-e"
//   - Lowercases and replaces remaining non-alphanumeric runs with hyphens
func toMetacriticSlug(title string) string {
	title = toASCII(title)
	var stripped strings.Builder
	for _, r := range title {
		if r == '\'' || r == '\u2019' || r == '.' {
			// drop: apostrophes merge adjacent words, dots merge abbreviations
			continue
		}
		stripped.WriteRune(r)
	}
	title = strings.ToLower(stripped.String())
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

// primaryTitle returns the portion of a title before the first ": " or " - ".
// e.g. "Jotun: Valhalla Edition" → "Jotun", "Hades - DLC" → "Hades"
func primaryTitle(title string) string {
	for _, sep := range []string{": ", " - "} {
		if idx := strings.Index(title, sep); idx >= 0 {
			return title[:idx]
		}
	}
	return title
}

// slugCandidates returns the ordered set of expected slug variants to try.
// Full title slug is tried first, then the primary title (before any subtitle).
func slugCandidates(title string) []string {
	seen := map[string]bool{}
	var result []string
	for _, t := range []string{title, primaryTitle(title)} {
		s := toMetacriticSlug(t)
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// bigrams returns the multiset of character bigrams in s.
func bigrams(s string) map[string]int {
	rs := []rune(s)
	bg := map[string]int{}
	for i := 0; i+1 < len(rs); i++ {
		bg[string(rs[i:i+2])]++
	}
	return bg
}

// bigramSim computes Sørensen-Dice bigram similarity between two strings.
func bigramSim(a, b string) float64 {
	if len([]rune(a)) < 2 || len([]rune(b)) < 2 {
		if a == b {
			return 1.0
		}
		return 0.0
	}
	ba, bb := bigrams(a), bigrams(b)
	intersection := 0
	for k, ca := range ba {
		if cb, ok := bb[k]; ok {
			if ca < cb {
				intersection += ca
			} else {
				intersection += cb
			}
		}
	}
	total := 0
	for _, c := range ba {
		total += c
	}
	for _, c := range bb {
		total += c
	}
	if total == 0 {
		return 0
	}
	return float64(2*intersection) / float64(total)
}

// mongeElkan computes Monge-Elkan similarity: for each query token find the
// max bigram similarity to any candidate token, then average those maxes.
func mongeElkan(queryTokens, candidateTokens []string) float64 {
	if len(queryTokens) == 0 || len(candidateTokens) == 0 {
		return 0
	}
	sum := 0.0
	for _, qt := range queryTokens {
		best := 0.0
		for _, ct := range candidateTokens {
			if s := bigramSim(qt, ct); s > best {
				best = s
			}
		}
		sum += best
	}
	return sum / float64(len(queryTokens))
}

// bestSlugMatch selects the best slug from allSlugs for the given title using:
//  1. Exact match against any expected slug variant
//  2. Prefix match (candidate starts with expectedSlug+"-")
//  3. Monge-Elkan fuzzy scoring against the full title slug tokens
func bestSlugMatch(expectedSlugs []string, allSlugs []string) string {
	// Exact match
	for _, slug := range allSlugs {
		for _, expected := range expectedSlugs {
			if slug == expected {
				return slug
			}
		}
	}
	// Prefix match
	for _, slug := range allSlugs {
		for _, expected := range expectedSlugs {
			if strings.HasPrefix(slug, expected+"-") {
				return slug
			}
		}
	}
	// Monge-Elkan fuzzy fallback using full title slug tokens
	if len(expectedSlugs) == 0 || len(allSlugs) == 0 {
		return ""
	}
	queryTokens := strings.Split(expectedSlugs[0], "-")
	bestScore := -1.0
	bestSlug := ""
	for _, slug := range allSlugs {
		candidateTokens := strings.Split(slug, "-")
		score := mongeElkan(queryTokens, candidateTokens)
		if score > bestScore {
			bestScore = score
			bestSlug = slug
		}
	}
	fmt.Printf("Fuzzy match: %q → %q (score=%.3f)\n", expectedSlugs[0], bestSlug, bestScore)
	return bestSlug
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
//  1. Fetch the search page and collect up to 50 /game/ slugs.
//  2. Pick the best slug via exact → prefix → Monge-Elkan fuzzy match.
//     Multiple expected slug variants are generated (full title + primary title).
//  3. Fetch the game's own page (/game/{slug}/) to extract the score.
func GetMetacriticScore(gameTitle string) (score int, metacriticSlug string, err error) {
	searchURL := "https://www.metacritic.com/search/" + url.PathEscape(gameTitle) + "/?category=13"
	fmt.Println("Search URL:", searchURL)

	body, statusCode, err := fetchMetacriticPage(searchURL)
	if err != nil {
		fmt.Println("Error fetching search page:", err)
		return 0, "", err
	}
	fmt.Println("Search response code:", statusCode)

	rawSlugs := reMetacriticSlug.FindAllStringSubmatch(body, 50)
	allSlugs := make([]string, 0, len(rawSlugs))
	seen := map[string]bool{}
	for _, m := range rawSlugs {
		if !seen[m[1]] {
			seen[m[1]] = true
			allSlugs = append(allSlugs, m[1])
		}
	}

	expected := slugCandidates(gameTitle)
	fmt.Printf("Looking for slug matching: %v\n", expected)

	metacriticSlug = bestSlugMatch(expected, allSlugs)

	if metacriticSlug == "" {
		fmt.Println("No game slugs found in search results")
		return 0, "", nil
	}

	// Fetch the game's own page for the score.
	gameURL := "https://www.metacritic.com/game/" + metacriticSlug + "/"
	fmt.Println("Game URL:", gameURL)

	gameBody, gameStatus, err := fetchMetacriticPage(gameURL)
	if err != nil {
		fmt.Println("Error fetching game page:", err)
		return 0, metacriticSlug, err
	}
	fmt.Println("Game page response code:", gameStatus)

	// If the critic score is TBD, the page uses global-score-tbd instead of global-score-value.
	// Bail out early to avoid accidentally matching a perfect individual review score.
	if strings.Contains(gameBody, `data-testid="global-score-tbd"`) {
		fmt.Println("Metascore is TBD on game page")
		return 0, metacriticSlug, nil
	}

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
