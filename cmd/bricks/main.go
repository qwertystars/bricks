package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/vadapavmov/bricks/internal/app"
)

const Version = "v1.3.1"
const MaxAllowedParallelDownloads = 5

func main() {
	// Define default values and usage messages for flags
	downloadUrl := flag.String("url", "", "Download url")
	downloadPath := flag.String("path", ".", "Download path")
	parallelDownloads := flag.Int("n", 3, "Number of parallel file downloads")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("Bricks %s\n", Version)
		return
	}

	// Get dirId
	baseURL, dirId := parseURL()

	// Build absolute path
	abspath, err := filepath.Abs(*downloadPath)
	if err != nil {
		log.Fatalf("invalid path %s", *downloadPath)
	}

	// Check if the specified download path exists
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		log.Fatalf("download path %s does not exist", abspath)
	}

	// To save site from DDOS
	if *parallelDownloads > MaxAllowedParallelDownloads {
		log.Fatalf("max parallel downloads can't be larger than %d", MaxAllowedParallelDownloads)
	}

	// Run the app
	bricks := app.New(baseURL)
	if err = bricks.Run(dirId, abspath, *parallelDownloads); err != nil {
		log.Fatalf("failed to download %v", err)
	}
}

func parseURL(downloadUrl string) (string, string) {
	var inputURL string
	if downloadUrl == "" {
		// Get the URL from stdin if not provided as a flag
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter URL: ")
		var err error
		inputURL, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read URL: %v", err)
		}
		inputURL = strings.TrimSpace(inputURL) // Trim newline and whitespaces
	} else {
		inputURL = strings.TrimSpace(downloadUrl)
	}
	
	// Parse URL and extract UUID
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}

	uuidStr := strings.TrimPrefix(parsedURL.Path, "/")
	uuidStr = strings.TrimSuffix(uuidStr, "/")

	// Validate UUID
	if _, err = uuid.Parse(uuidStr); err != nil {
		log.Fatalf("invalid UUID format in URL: %v", err)
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	return baseURL, uuidStr
}
