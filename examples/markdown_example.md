# SIG Pack to Markdown Conversion

The `sigma` tool now includes a `markdown` command that converts SIQ files to readable markdown format.

## Usage

```bash
go run sigma.go markdown [siq-file] [output-file]
```

## Example

```bash
# Convert a SIQ file to markdown
go run sigma.go markdown data/imbaenergy.siq output.md
```

## Output Format

The generated markdown file includes only the essential content:

- **Rounds**: Each round with its name
- **Themes**: Each theme within rounds
- **Questions**: Question content and right answers
- **Right Answers**: Single answer (no numbering) or numbered list for multiple answers

## Features

- Supports both SIQ v4 and v5 formats
- Clean, minimal output without metadata
- Simple formatting for single answers
- Numbered lists only for multiple correct answers
- Preserves question content (text, images, audio, video)

## Sample Output

```markdown
## Round 1: Round Name

### Theme 1: Theme Name

#### Question 1

**Content**:
- **image**: @image.png

**Right Answer**:
Correct Answer

---

#### Question 2

**Content**:
- **text**: Question text with multiple choice options

**Right Answers**:
1. First correct answer
2. Second correct answer
``` 