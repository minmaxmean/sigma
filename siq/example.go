package siq

import (
	"fmt"
	"log"
)

// Example demonstrates how to use the SIQ library
func Example() {
	// Open a SIQ file
	reader, err := NewSIQReader("example.siq")
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer reader.Close()

	// Read and parse the package
	pkg, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read SIQ file:", err)
	}

	// Print package information
	fmt.Printf("Package: %s\n", pkg.Name)
	fmt.Printf("Version: %s\n", pkg.Version)
	fmt.Printf("Difficulty: %d/10\n", pkg.Difficulty)
	fmt.Printf("Language: %s\n", pkg.Language)
	fmt.Printf("Publisher: %s\n", pkg.Publisher)
	fmt.Printf("Date: %s\n", pkg.Date)

	// Print statistics
	fmt.Printf("\nStatistics:\n")
	fmt.Printf("Rounds: %d\n", pkg.GetRoundCount())
	fmt.Printf("Themes: %d\n", pkg.GetThemeCount())
	fmt.Printf("Questions: %d\n", pkg.GetQuestionCount())

	// List all files in the archive
	fmt.Printf("\nFiles in archive:\n")
	for _, fileName := range reader.ListFiles() {
		fmt.Printf("  %s\n", fileName)
	}

	// Process all questions
	fmt.Printf("\nQuestions:\n")
	questions := pkg.GetAllQuestions()
	for i, question := range questions {
		fmt.Printf("Question %d:\n", i+1)
		fmt.Printf("  Type: %s\n", question.Type)
		fmt.Printf("  Right answers: %d\n", len(question.Right))
		fmt.Printf("  Wrong answers: %d\n", len(question.Wrong))

		// Get question content
		content := question.GetQuestionContent()
		if len(content) > 0 {
			fmt.Printf("  Content items: %d\n", len(content))
			for j, item := range content {
				fmt.Printf("    Item %d: %s (%s)\n", j+1, item.Type, item.Value)
			}
		}
		fmt.Println()
	}
}

// ExampleExtractFiles demonstrates how to extract files from the SIQ archive
func ExampleExtractFiles() {
	reader, err := NewSIQReader("example.siq")
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer reader.Close()

	// Extract all images
	files := reader.ListFiles()
	for _, fileName := range files {
		if isImageFile(fileName) {
			destPath := fmt.Sprintf("extracted/%s", fileName)
			if err := reader.ExtractFile(fileName, destPath); err != nil {
				log.Printf("Failed to extract %s: %v", fileName, err)
			} else {
				fmt.Printf("Extracted: %s\n", fileName)
			}
		}
	}
}

// ExampleResolveReferences demonstrates how to resolve references
func ExampleResolveReferences() {
	reader, err := NewSIQReader("example.siq")
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer reader.Close()

	pkg, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read SIQ file:", err)
	}

	// Example of resolving a reference
	ref := "@ae9f7eb2-6091-4b34-97a1-0f74ad193d57"
	resolved, err := pkg.ResolveReference(ref)
	if err != nil {
		fmt.Printf("Failed to resolve reference %s: %v\n", ref, err)
	} else {
		fmt.Printf("Resolved %s to: %s\n", ref, resolved)
	}
}

// Helper function to check if a file is an image
func isImageFile(fileName string) bool {
	fileExt := getFileExtension(fileName)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	for _, ext := range imageExts {
		if fileExt == ext {
			return true
		}
	}
	return false
}

// Helper function to get file extension
func getFileExtension(fileName string) string {
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ""
}
