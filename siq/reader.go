package siq

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Package represents the main SIG pack structure (supports both v4 and v5)
type Package struct {
	ID          string  `xml:"id,attr"`
	Name        string  `xml:"name,attr"`
	Version     string  `xml:"version,attr"`
	Restriction string  `xml:"restriction,attr"`
	Date        string  `xml:"date,attr"`
	Publisher   string  `xml:"publisher,attr"`
	Difficulty  int     `xml:"difficulty,attr"`
	Logo        string  `xml:"logo,attr"`
	Language    string  `xml:"language,attr"`
	Info        *Info   `xml:"info,omitempty"`
	Tags        *Tags   `xml:"tags,omitempty"`
	Global      *Global `xml:"global,omitempty"`
	Rounds      []Round `xml:"round"`
	// Version 4 compatibility
	RoundsV4 []RoundV4 `xml:"rounds>round,omitempty"`
}

// Round represents a round in the SIG pack (v5 format)
type Round struct {
	Name   string  `xml:"name,attr"`
	Info   *Info   `xml:"info,omitempty"`
	Themes []Theme `xml:"theme"`
}

// RoundV4 represents a round in the SIG pack (v4 format)
type RoundV4 struct {
	Name   string    `xml:"name,attr"`
	Info   *Info     `xml:"info,omitempty"`
	Themes []ThemeV4 `xml:"themes>theme"`
}

// Theme represents a theme in a round (v5 format)
type Theme struct {
	Name      string     `xml:"name,attr"`
	Info      *Info      `xml:"info,omitempty"`
	Questions []Question `xml:"question"`
}

// ThemeV4 represents a theme in a round (v4 format)
type ThemeV4 struct {
	Name      string       `xml:"name,attr"`
	Info      *Info        `xml:"info,omitempty"`
	Questions []QuestionV4 `xml:"questions>question"`
}

// Question represents a question in a theme (v5 format)
type Question struct {
	Type   string   `xml:"type,attr"`
	Params []Param  `xml:"params>param"`
	Right  []string `xml:"right>answer"`
	Wrong  []string `xml:"wrong>answer"`
	Script *Script  `xml:"script,omitempty"`
	Info   *Info    `xml:"info,omitempty"`
}

// QuestionV4 represents a question in a theme (v4 format)
type QuestionV4 struct {
	Price    int       `xml:"price,attr"`
	Scenario *Scenario `xml:"scenario,omitempty"`
	Right    []string  `xml:"right>answer"`
	Wrong    []string  `xml:"wrong>answer"`
	Info     *Info     `xml:"info,omitempty"`
}

// Scenario represents a scenario in v4 format
type Scenario struct {
	Atoms []Atom `xml:"atom"`
}

// Atom represents an atom in v4 format
type Atom struct {
	Type     string `xml:"type,attr"`
	Duration int    `xml:"duration,attr,omitempty"`
	Content  string `xml:",chardata"`
}

// Param represents a parameter in a question (v5 format)
type Param struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
	// For content type parameters
	Items []ContentItem `xml:"item,omitempty"`
	// For group type parameters
	Params []Param `xml:"param,omitempty"`
	// For numberSet type parameters
	Minimum int `xml:"minimum,omitempty"`
	Maximum int `xml:"maximum,omitempty"`
	Step    int `xml:"step,omitempty"`
}

// ContentItem represents a content item in a question (v5 format)
type ContentItem struct {
	Type          string `xml:"type,attr"`
	Value         string `xml:",chardata"`
	Duration      int    `xml:"duration,attr,omitempty"`
	IsRef         bool   `xml:"isRef,attr,omitempty"`
	Placement     string `xml:"placement,attr,omitempty"`
	WaitForFinish bool   `xml:"waitForFinish,attr,omitempty"`
}

// Script represents a script for complex scenarios
type Script struct {
	Content string `xml:",chardata"`
}

// Info represents information about authors, sources, comments, etc.
type Info struct {
	Authors         []string `xml:"authors>author"`
	Sources         []string `xml:"sources>source"`
	Comments        []string `xml:"comments>comment"`
	ShowmanComments []string `xml:"showmanComments>comment"`
}

// Tags represents package tags
type Tags struct {
	Tags []string `xml:"tag"`
}

// Global represents globally defined authors and sources
type Global struct {
	Authors []GlobalAuthor `xml:"authors>author"`
	Sources []GlobalSource `xml:"sources>source"`
}

// GlobalAuthor represents a globally defined author
type GlobalAuthor struct {
	ID   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

// GlobalSource represents a globally defined source
type GlobalSource struct {
	ID   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

// SIQReader represents a reader for SIQ files
type SIQReader struct {
	zipReader *zip.ReadCloser
	pkg       *Package
	global    *Global
	version   int // 4 or 5
}

// NewSIQReader creates a new SIQ reader
func NewSIQReader(filePath string) (*SIQReader, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SIQ file: %w", err)
	}

	return &SIQReader{
		zipReader: zipReader,
	}, nil
}

// Read reads and parses the SIQ file
func (r *SIQReader) Read() (*Package, error) {
	// Find and read content.xml
	contentFile, err := r.findContentFile()
	if err != nil {
		return nil, err
	}

	// Parse the XML content
	if err := r.parseContent(contentFile); err != nil {
		return nil, err
	}

	return r.pkg, nil
}

// findContentFile finds the content.xml file in the zip archive
func (r *SIQReader) findContentFile() (*zip.File, error) {
	for _, file := range r.zipReader.File {
		if file.Name == "content.xml" {
			return file, nil
		}
	}
	return nil, fmt.Errorf("content.xml not found in SIQ file")
}

// parseContent parses the content.xml file
func (r *SIQReader) parseContent(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open content.xml: %w", err)
	}
	defer rc.Close()

	// First, try to detect the version by reading the XML declaration
	content, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read content.xml: %w", err)
	}

	// Check if it's version 4 (has xmlns="http://vladimirkhil.com/ygpackage3.0.xsd")
	if strings.Contains(string(content), `xmlns="http://vladimirkhil.com/ygpackage3.0.xsd"`) {
		r.version = 4
	} else {
		r.version = 5
	}

	// Parse based on version
	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	if r.version == 4 {
		if err := decoder.Decode(&r.pkg); err != nil {
			return fmt.Errorf("failed to decode v4 XML: %w", err)
		}
	} else {
		if err := decoder.Decode(&r.pkg); err != nil {
			return fmt.Errorf("failed to decode v5 XML: %w", err)
		}
	}

	return nil
}

// Close closes the SIQ reader
func (r *SIQReader) Close() error {
	if r.zipReader != nil {
		return r.zipReader.Close()
	}
	return nil
}

// GetFile retrieves a file from the SIQ archive
func (r *SIQReader) GetFile(filePath string) (*zip.File, error) {
	// Handle URI-encoded file names
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		decodedPath = filePath
	}

	for _, file := range r.zipReader.File {
		if file.Name == decodedPath || file.Name == filePath {
			return file, nil
		}
	}
	return nil, fmt.Errorf("file %s not found in SIQ archive", filePath)
}

// ListFiles lists all files in the SIQ archive
func (r *SIQReader) ListFiles() []string {
	var files []string
	for _, file := range r.zipReader.File {
		files = append(files, file.Name)
	}
	return files
}

// ExtractFile extracts a file from the SIQ archive to a destination path
func (r *SIQReader) ExtractFile(filePath, destPath string) error {
	file, err := r.GetFile(filePath)
	if err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer rc.Close()

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, rc)
	if err != nil {
		return fmt.Errorf("failed to copy file %s: %w", filePath, err)
	}

	return nil
}

// GetVersion returns the detected SIQ format version
func (r *SIQReader) GetVersion() int {
	return r.version
}

// GetQuestionCount returns the total number of questions in the package
func (p *Package) GetQuestionCount() int {
	count := 0

	// Count v5 questions
	for _, round := range p.Rounds {
		for _, theme := range round.Themes {
			count += len(theme.Questions)
		}
	}

	// Count v4 questions
	for _, round := range p.RoundsV4 {
		for _, theme := range round.Themes {
			count += len(theme.Questions)
		}
	}

	return count
}

// GetRoundCount returns the number of rounds in the package
func (p *Package) GetRoundCount() int {
	return len(p.Rounds) + len(p.RoundsV4)
}

// GetThemeCount returns the total number of themes in the package
func (p *Package) GetThemeCount() int {
	count := 0

	// Count v5 themes
	for _, round := range p.Rounds {
		count += len(round.Themes)
	}

	// Count v4 themes
	for _, round := range p.RoundsV4 {
		count += len(round.Themes)
	}

	return count
}

// GetAllQuestions returns all questions from the package (converted to v5 format)
func (p *Package) GetAllQuestions() []Question {
	var questions []Question

	// Add v5 questions
	for _, round := range p.Rounds {
		for _, theme := range round.Themes {
			questions = append(questions, theme.Questions...)
		}
	}

	// Convert v4 questions to v5 format
	for _, round := range p.RoundsV4 {
		for _, theme := range round.Themes {
			for _, qv4 := range theme.Questions {
				question := convertV4ToV5Question(qv4)
				questions = append(questions, question)
			}
		}
	}

	return questions
}

// GetQuestionsByType returns all questions of a specific type
func (p *Package) GetQuestionsByType(questionType string) []Question {
	var questions []Question
	allQuestions := p.GetAllQuestions()

	for _, question := range allQuestions {
		if question.Type == questionType {
			questions = append(questions, question)
		}
	}

	return questions
}

// ConvertV4ToV5Question converts a v4 question to v5 format (exported version)
func ConvertV4ToV5Question(qv4 QuestionV4) Question {
	return convertV4ToV5Question(qv4)
}

// convertV4ToV5Question converts a v4 question to v5 format
func convertV4ToV5Question(qv4 QuestionV4) Question {
	question := Question{
		Type:  "simple", // Default type for v4 questions
		Right: qv4.Right,
		Wrong: qv4.Wrong,
		Info:  qv4.Info,
	}

	// Convert scenario atoms to content items
	if qv4.Scenario != nil {
		var items []ContentItem
		for _, atom := range qv4.Scenario.Atoms {
			item := ContentItem{
				Type:     atom.Type,
				Value:    atom.Content,
				Duration: atom.Duration,
			}
			items = append(items, item)
		}

		// Create a question parameter with the content items
		param := Param{
			Name:  "question",
			Type:  "content",
			Items: items,
		}
		question.Params = append(question.Params, param)
	}

	return question
}

// GetQuestionContent returns the content items for a question
func (q *Question) GetQuestionContent() []ContentItem {
	for _, param := range q.Params {
		if param.Name == "question" {
			return param.Items
		}
	}
	return nil
}

// GetParamValue returns the value of a parameter by name
func (q *Question) GetParamValue(name string) string {
	for _, param := range q.Params {
		if param.Name == name {
			return param.Value
		}
	}
	return ""
}

// GetParamItems returns the content items of a parameter by name
func (q *Question) GetParamItems(name string) []ContentItem {
	for _, param := range q.Params {
		if param.Name == name {
			return param.Items
		}
	}
	return nil
}

// ResolveReference resolves a reference starting with @
func (p *Package) ResolveReference(ref string) (string, error) {
	if !strings.HasPrefix(ref, "@") {
		return ref, nil
	}

	// Extract ID and specification
	parts := strings.SplitN(ref[1:], "#", 2)
	id := parts[0]
	specification := ""
	if len(parts) > 1 {
		specification = parts[1]
	}

	// Look up in global authors/sources
	if p.Global != nil {
		for _, author := range p.Global.Authors {
			if author.ID == id {
				if specification != "" {
					return author.Name + " " + specification, nil
				}
				return author.Name, nil
			}
		}
		for _, source := range p.Global.Sources {
			if source.ID == id {
				if specification != "" {
					return source.Name + " " + specification, nil
				}
				return source.Name, nil
			}
		}
	}

	return ref, fmt.Errorf("reference %s not found", ref)
}
