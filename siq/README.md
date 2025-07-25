# SIQ Library

A Go library for reading and parsing SIG pack files (.siq format) according to the SIQ file format version 5 specification.

## Features

- **Complete SIQ Format Support**: Full support for SIQ file format version 5
- **ZIP Archive Handling**: Automatic handling of SIQ files as ZIP archives
- **XML Parsing**: Robust XML parsing with proper structure mapping
- **Reference Resolution**: Support for @ references to global authors and sources
- **File Extraction**: Extract multimedia files from SIQ archives
- **URI Encoding Support**: Handle URI-encoded file names for backward compatibility
- **Comprehensive API**: Easy-to-use API for accessing package data

## Installation

```bash
go get github.com/minmaxmean/sigma/siq
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/minmaxmean/sigma/siq"
)

func main() {
    // Open a SIQ file
    reader, err := siq.NewSIQReader("example.siq")
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
    fmt.Printf("Difficulty: %d/10\n", pkg.Difficulty)
    fmt.Printf("Questions: %d\n", pkg.GetQuestionCount())
}
```

## API Reference

### Main Types

#### Package
Represents the main SIG pack structure.

```go
type Package struct {
    ID           string   `xml:"id,attr"`
    Name         string   `xml:"name,attr"`
    Version      string   `xml:"version,attr"`
    Restriction  string   `xml:"restriction,attr"`
    Date         string   `xml:"date,attr"`
    Publisher    string   `xml:"publisher,attr"`
    Difficulty   int      `xml:"difficulty,attr"`
    Logo         string   `xml:"logo,attr"`
    Language     string   `xml:"language,attr"`
    Info         *Info    `xml:"info,omitempty"`
    Tags         *Tags    `xml:"tags,omitempty"`
    Global       *Global  `xml:"global,omitempty"`
    Rounds       []Round  `xml:"round"`
}
```

#### Question
Represents a question in a theme.

```go
type Question struct {
    Type   string   `xml:"type,attr"`
    Params []Param  `xml:"params>param"`
    Right  []string `xml:"right>answer"`
    Wrong  []string `xml:"wrong>answer"`
    Script *Script  `xml:"script,omitempty"`
    Info   *Info    `xml:"info,omitempty"`
}
```

#### ContentItem
Represents a content item in a question.

```go
type ContentItem struct {
    Type         string `xml:"type,attr"`
    Value        string `xml:",chardata"`
    Duration     int    `xml:"duration,attr,omitempty"`
    IsRef        bool   `xml:"isRef,attr,omitempty"`
    Placement    string `xml:"placement,attr,omitempty"`
    WaitForFinish bool  `xml:"waitForFinish,attr,omitempty"`
}
```

### Main Functions

#### NewSIQReader
Creates a new SIQ reader for a file.

```go
reader, err := siq.NewSIQReader("file.siq")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()
```

#### Read
Reads and parses the SIQ file.

```go
pkg, err := reader.Read()
if err != nil {
    log.Fatal(err)
}
```

### Package Methods

#### GetQuestionCount
Returns the total number of questions in the package.

```go
count := pkg.GetQuestionCount()
```

#### GetAllQuestions
Returns all questions from the package.

```go
questions := pkg.GetAllQuestions()
```

#### GetQuestionsByType
Returns all questions of a specific type.

```go
simpleQuestions := pkg.GetQuestionsByType("simple")
```

#### ResolveReference
Resolves a reference starting with @.

```go
resolved, err := pkg.ResolveReference("@author-id")
if err != nil {
    log.Fatal(err)
}
```

### Question Methods

#### GetQuestionContent
Returns the content items for a question.

```go
content := question.GetQuestionContent()
for _, item := range content {
    fmt.Printf("Type: %s, Value: %s\n", item.Type, item.Value)
}
```

#### GetParamValue
Returns the value of a parameter by name.

```go
value := question.GetParamValue("param-name")
```

#### GetParamItems
Returns the content items of a parameter by name.

```go
items := question.GetParamItems("param-name")
```

### SIQReader Methods

#### ListFiles
Lists all files in the SIQ archive.

```go
files := reader.ListFiles()
for _, file := range files {
    fmt.Println(file)
}
```

#### GetFile
Retrieves a file from the SIQ archive.

```go
file, err := reader.GetFile("Images/logo.png")
if err != nil {
    log.Fatal(err)
}
```

#### ExtractFile
Extracts a file from the SIQ archive to a destination path.

```go
err := reader.ExtractFile("Images/logo.png", "extracted/logo.png")
if err != nil {
    log.Fatal(err)
}
```

## Question Types

The library supports all well-known question types:

- `simple` - Simple question
- `cat` - Cat in the bag
- `auction` - Auction
- `bagCat` - Bag cat
- `spider` - Spider
- `secret` - Secret
- `noRisk` - No risk
- `super` - Super
- `complex` - Complex
- `media` - Media
- `stake` - Stake
- `final` - Final

## Content Types

Supported content item types:

- `text` - Text content
- `image` - Image content
- `audio` - Audio content
- `video` - Video content
- `html` - HTML content

## Placement Types

Content placement options:

- `screen` - Displayed on game screen (default)
- `replic` - Provided as showman replic
- `background` - Played in the background

## Parameter Types

Parameter type options:

- `simple` - Simple string value (default)
- `content` - Contains a set of content items
- `group` - Contains a set of other parameters
- `numberSet` - Represents a range of numbers

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/minmaxmean/sigma/siq"
)

func main() {
    reader, err := siq.NewSIQReader("game.siq")
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    pkg, err := reader.Read()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Package: %s\n", pkg.Name)
    fmt.Printf("Questions: %d\n", pkg.GetQuestionCount())

    // Process all questions
    for i, question := range pkg.GetAllQuestions() {
        fmt.Printf("Question %d: %s\n", i+1, question.Type)
        
        content := question.GetQuestionContent()
        for _, item := range content {
            fmt.Printf("  Content: %s (%s)\n", item.Type, item.Value)
        }
    }
}
```

### Extract Multimedia Files

```go
func extractMedia(reader *siq.SIQReader) {
    files := reader.ListFiles()
    for _, fileName := range files {
        if strings.HasPrefix(fileName, "Images/") ||
           strings.HasPrefix(fileName, "Audio/") ||
           strings.HasPrefix(fileName, "Video/") {
            destPath := fmt.Sprintf("extracted/%s", fileName)
            if err := reader.ExtractFile(fileName, destPath); err != nil {
                log.Printf("Failed to extract %s: %v", fileName, err)
            } else {
                fmt.Printf("Extracted: %s\n", fileName)
            }
        }
    }
}
```

### Resolve References

```go
func resolveReferences(pkg *siq.Package) {
    // Resolve author references
    for _, round := range pkg.Rounds {
        for _, theme := range round.Themes {
            if theme.Info != nil {
                for _, author := range theme.Info.Authors {
                    resolved, err := pkg.ResolveReference(author)
                    if err != nil {
                        fmt.Printf("Failed to resolve %s: %v\n", author, err)
                    } else {
                        fmt.Printf("Author: %s\n", resolved)
                    }
                }
            }
        }
    }
}
```

## Testing

Run the tests:

```bash
go test ./siq
```

## License

This library is part of the sigma project and follows the same license terms.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 