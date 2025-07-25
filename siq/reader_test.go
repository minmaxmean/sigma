package siq

import (
	"archive/zip"
	"os"
	"testing"
)

// createTestSIQFile creates a test SIQ file for testing
func createTestSIQFile(t *testing.T) string {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.siq")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer tmpFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	// Create sample content.xml
	contentXML := `<?xml version="1.0" encoding="UTF-8"?>
<package id="test-package" name="Test Package" version="5.0" difficulty="5" language="en">
	<info>
		<authors>
			<author>Test Author</author>
		</authors>
		<sources>
			<source>Test Source</source>
		</sources>
	</info>
	<round name="Round 1">
		<theme name="Test Theme">
			<question type="simple">
				<params>
					<param name="question" type="content">
						<item type="text">What is 2+2?</item>
					</param>
				</params>
				<right>
					<answer>4</answer>
				</right>
				<wrong>
					<answer>3</answer>
					<answer>5</answer>
				</wrong>
			</question>
		</theme>
	</round>
</package>`

	// Add content.xml to the zip
	contentFile, err := zipWriter.Create("content.xml")
	if err != nil {
		t.Fatal("Failed to create content.xml in zip:", err)
	}
	_, err = contentFile.Write([]byte(contentXML))
	if err != nil {
		t.Fatal("Failed to write content.xml:", err)
	}

	return tmpFile.Name()
}

func TestNewSIQReader(t *testing.T) {
	// Create a test file
	testFile := createTestSIQFile(t)
	defer os.Remove(testFile)

	// Test creating a new reader
	reader, err := NewSIQReader(testFile)
	if err != nil {
		t.Fatal("Failed to create SIQ reader:", err)
	}
	defer reader.Close()

	if reader.zipReader == nil {
		t.Error("zipReader should not be nil")
	}
}

func TestRead(t *testing.T) {
	// Create a test file
	testFile := createTestSIQFile(t)
	defer os.Remove(testFile)

	// Create reader and read package
	reader, err := NewSIQReader(testFile)
	if err != nil {
		t.Fatal("Failed to create SIQ reader:", err)
	}
	defer reader.Close()

	pkg, err := reader.Read()
	if err != nil {
		t.Fatal("Failed to read SIQ file:", err)
	}

	// Verify package data
	if pkg.ID != "test-package" {
		t.Errorf("Expected package ID 'test-package', got '%s'", pkg.ID)
	}
	if pkg.Name != "Test Package" {
		t.Errorf("Expected package name 'Test Package', got '%s'", pkg.Name)
	}
	if pkg.Version != "5.0" {
		t.Errorf("Expected version '5.0', got '%s'", pkg.Version)
	}
	if pkg.Difficulty != 5 {
		t.Errorf("Expected difficulty 5, got %d", pkg.Difficulty)
	}
	if pkg.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", pkg.Language)
	}

	// Verify structure
	if len(pkg.Rounds) != 1 {
		t.Errorf("Expected 1 round, got %d", len(pkg.Rounds))
	}
	if len(pkg.Rounds[0].Themes) != 1 {
		t.Errorf("Expected 1 theme, got %d", len(pkg.Rounds[0].Themes))
	}
	if len(pkg.Rounds[0].Themes[0].Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(pkg.Rounds[0].Themes[0].Questions))
	}

	// Verify question
	question := pkg.Rounds[0].Themes[0].Questions[0]
	if question.Type != "simple" {
		t.Errorf("Expected question type 'simple', got '%s'", question.Type)
	}
	if len(question.Right) != 1 {
		t.Errorf("Expected 1 right answer, got %d", len(question.Right))
	}
	if len(question.Wrong) != 2 {
		t.Errorf("Expected 2 wrong answers, got %d", len(question.Wrong))
	}
}

func TestGetQuestionContent(t *testing.T) {
	// Create a test file
	testFile := createTestSIQFile(t)
	defer os.Remove(testFile)

	// Create reader and read package
	reader, err := NewSIQReader(testFile)
	if err != nil {
		t.Fatal("Failed to create SIQ reader:", err)
	}
	defer reader.Close()

	pkg, err := reader.Read()
	if err != nil {
		t.Fatal("Failed to read SIQ file:", err)
	}

	// Get the question
	question := pkg.Rounds[0].Themes[0].Questions[0]
	content := question.GetQuestionContent()

	if len(content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(content))
	}

	if content[0].Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", content[0].Type)
	}
	if content[0].Value != "What is 2+2?" {
		t.Errorf("Expected content value 'What is 2+2?', got '%s'", content[0].Value)
	}
}

func TestGetQuestionCount(t *testing.T) {
	// Create a test file
	testFile := createTestSIQFile(t)
	defer os.Remove(testFile)

	// Create reader and read package
	reader, err := NewSIQReader(testFile)
	if err != nil {
		t.Fatal("Failed to create SIQ reader:", err)
	}
	defer reader.Close()

	pkg, err := reader.Read()
	if err != nil {
		t.Fatal("Failed to read SIQ file:", err)
	}

	count := pkg.GetQuestionCount()
	if count != 1 {
		t.Errorf("Expected 1 question, got %d", count)
	}
}

func TestListFiles(t *testing.T) {
	// Create a test file
	testFile := createTestSIQFile(t)
	defer os.Remove(testFile)

	// Create reader
	reader, err := NewSIQReader(testFile)
	if err != nil {
		t.Fatal("Failed to create SIQ reader:", err)
	}
	defer reader.Close()

	files := reader.ListFiles()
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
	if files[0] != "content.xml" {
		t.Errorf("Expected 'content.xml', got '%s'", files[0])
	}
}

func TestResolveReference(t *testing.T) {
	// Create a test package with global references
	pkg := &Package{
		Global: &Global{
			Authors: []GlobalAuthor{
				{ID: "author1", Name: "John Doe"},
			},
			Sources: []GlobalSource{
				{ID: "source1", Name: "Test Book"},
			},
		},
	}

	// Test resolving author reference
	resolved, err := pkg.ResolveReference("@author1")
	if err != nil {
		t.Errorf("Failed to resolve author reference: %v", err)
	}
	if resolved != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", resolved)
	}

	// Test resolving source reference with specification
	resolved, err = pkg.ResolveReference("@source1#p.123")
	if err != nil {
		t.Errorf("Failed to resolve source reference: %v", err)
	}
	if resolved != "Test Book p.123" {
		t.Errorf("Expected 'Test Book p.123', got '%s'", resolved)
	}

	// Test non-reference
	resolved, err = pkg.ResolveReference("Direct text")
	if err != nil {
		t.Errorf("Failed to handle direct text: %v", err)
	}
	if resolved != "Direct text" {
		t.Errorf("Expected 'Direct text', got '%s'", resolved)
	}
}
