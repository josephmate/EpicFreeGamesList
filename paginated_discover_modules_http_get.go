package main

import (
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)


const (
	PaginatedDiscoverModulesAndroidURL = "https://egs-platform-service.store.epicgames.com/api/v2/public/discover/home?count=10&country=US&locale=en&platform=android&start=0&store=EGS"
	PaginatedDiscoverModulesIosURL     = "https://egs-platform-service.store.epicgames.com/api/v2/public/discover/home?count=10&country=US&locale=en&platform=ios&start=0&store=EGS"
)

// TLS client with timeout and Chrome profile
var tlsClient tls_client.HttpClient

func init() {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_107),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		panic(fmt.Sprintf("failed to create TLS client: %v", err))
	}
	tlsClient = client
}

// HttpGetPaginatedDiscoverModules makes an HTTP GET request using TLS client and returns the response body as a string
func HttpGetPaginatedDiscoverModules(url string) (string, error) {
	// Create HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a real browser request
	req.Header = http.Header{
		"accept":             {"*/*"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
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

	// Make the request
	resp, err := tlsClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}