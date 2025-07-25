package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug/main.go <siq-file>")
		os.Exit(1)
	}

	siqFile := os.Args[1]

	// Open the zip file
	zipReader, err := zip.OpenReader(siqFile)
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer zipReader.Close()

	// Find content.xml
	var contentFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "content.xml" {
			contentFile = file
			break
		}
	}

	if contentFile == nil {
		log.Fatal("content.xml not found in SIQ file")
	}

	// Read and display the first 1000 characters of content.xml
	rc, err := contentFile.Open()
	if err != nil {
		log.Fatal("Failed to open content.xml:", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal("Failed to read content.xml:", err)
	}

	fmt.Printf("=== Content.xml Structure ===\n")
	fmt.Printf("File size: %d bytes\n\n", len(content))

	// Show first 2000 characters
	displayLength := 2000
	if len(content) < displayLength {
		displayLength = len(content)
	}

	fmt.Printf("First %d characters:\n", displayLength)
	fmt.Printf("%s\n", string(content[:displayLength]))

	if len(content) > displayLength {
		fmt.Printf("\n... (truncated, total %d bytes)\n", len(content))
	}
}
