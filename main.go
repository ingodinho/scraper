package main

import (
	"fmt"
	"os"

	"github.com/ingodinho/crawler/internal"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	fmt.Printf("starting crawl of: %v\n", baseURL)

	pages := map[string]int {}
	err := internal.CrawlPage(baseURL, baseURL, pages)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	fmt.Println(pages)
	os.Exit(0)
}
