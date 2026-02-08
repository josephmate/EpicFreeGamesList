package main

import (
    "flag"
    "fmt"
    "os"
)

func CliHandlerMobileDiscoverPage() {
    fs := flag.NewFlagSet("mobile_discover_page", flag.ExitOnError)
    platform := fs.String("platform", "", "Platform to fetch data for (android|ios)")
    
    fs.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: tool mobile_discover_page --platform <android|ios>\n")
        fmt.Fprintf(os.Stderr, "\nOptions:\n")
        fs.PrintDefaults()
        fmt.Fprintf(os.Stderr, "\nExample:\n")
        fmt.Fprintf(os.Stderr, "  tool mobile_discover_page --platform android\n")
        fmt.Fprintf(os.Stderr, "  tool mobile_discover_page --platform ios\n")
    }

    // Parse command line arguments starting from index 2 (skip "tool" and "mobile_discover_page")
    fs.Parse(os.Args[2:])

    // Validate platform argument
    if len(*platform) == 0 {
        fmt.Fprintf(os.Stderr, "Error: --platform argument is required\n\n")
        fs.Usage()
        os.Exit(1)
    }

    // Determine URL based on platform
    var url string
    switch *platform {
    case "android":
        url = PaginatedDiscoverModulesAndroidURL
        fmt.Printf("Fetching Android discover modules...\n")
    case "ios":
        url = PaginatedDiscoverModulesIosURL
        fmt.Printf("Fetching iOS discover modules...\n")
    default:
        fmt.Fprintf(os.Stderr, "Error: invalid platform '%s'. Must be 'android' or 'ios'\n\n", *platform)
        fs.Usage()
        os.Exit(1)
    }

    // Fetch discover modules
    result, err := HttpGetPaginatedDiscoverModules(url)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error fetching discover modules: %v\n", err)
        os.Exit(1)
    }

    // Display results
    fmt.Printf("✅ Successfully fetched discover modules for %s!\n\n", *platform)

    // has non-empty string
    if len(result) > 0 {
        result, err := ParsePaginatedDiscoverModules(result)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error fetching discover modules: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("🎮 Found %d discover modules:\n", len(result.Data))
        
        for i, module := range result.Data {
            if module == nil {
                continue
            }

            fmt.Printf("\n%d. Module: %s\n", i+1, getStringValue(module.Title))
            fmt.Printf("   Type: %s\n", getStringValue(module.Type))
            fmt.Printf("   Size: %s\n", getStringValue(module.Size))
            fmt.Printf("   Topic ID: %s\n", getStringValue(module.TopicID))

            if len(module.Offers) > 0 {
                fmt.Printf("   Offers: %d\n", len(module.Offers))
                
                // Show details of first offer
                if module.Offers[0] != nil && module.Offers[0].Content != nil {
                    offer := module.Offers[0]
                    fmt.Printf("   First Offer: %s\n", getStringValue(offer.Content.Title))
                    
                    // Show purchase info if available
                    if len(offer.Content.Purchase) > 0 && offer.Content.Purchase[0] != nil {
                        purchase := offer.Content.Purchase[0]
                        if purchase.Price != nil {
                            fmt.Printf("   Price: %.2f %s\n", 
                                getFloat64Value(purchase.Price.DecimalPrice),
                                getStringValue(purchase.Price.CurrencyCode))
                        }
                        fmt.Printf("   Purchase Type: %s\n", getStringValue(purchase.PurchaseType))
                    }

                    // Show platform info if available
                    if offer.Content.SystemSpecs != nil {
                        fmt.Printf("   Platform: %s\n", getStringValue(offer.Content.SystemSpecs.Platform))
                        if offer.Content.SystemSpecs.ApplicationVersion != nil {
                            fmt.Printf("   Version: %s\n", getStringValue(offer.Content.SystemSpecs.ApplicationVersion))
                        }
                    }
                }
            } else {
                fmt.Printf("   Offers: 0\n")
            }

            if module.Link != nil {
                fmt.Printf("   Link: %s (%s)\n", getStringValue(module.Link.Src), getStringValue(module.Link.LinkText))
            }
        }
    } else {
        fmt.Printf("No discover modules found.\n")
    }

    fmt.Printf("\n🎉 Done!\n")
}


// Helper functions to safely get values from pointers
func getStringValue(s *string) string {
    if s == nil {
        return "N/A"
    }
    return *s
}

func getIntValue(i *int) int {
    if i == nil {
        return 0
    }
    return *i
}

func getFloat64Value(f *float64) float64 {
    if f == nil {
        return 0.0
    }
    return *f
}
