package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minmaxmean/sigma/siq"
	"github.com/spf13/cobra"
)

var markdownCmd = &cobra.Command{
	Use:   "markdown [siq-file] [output-file]",
	Short: "Convert a SIQ file to markdown format",
	Long: `Convert a SIQ file to a readable markdown format including:
- Package metadata and statistics
- All rounds, themes, and questions
- Question content and answers
- Authors, sources, and comments`,
	Args: cobra.ExactArgs(2),
	Run:  runMarkdown,
}

func runMarkdown(cmd *cobra.Command, args []string) {
	siqFile := args[0]
	outputFile := args[1]

	// Open and read the SIQ file
	reader, err := siq.NewSIQReader(siqFile)
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer reader.Close()

	// Parse the package
	pkg, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read SIQ file:", err)
	}

	// Generate markdown content
	markdown := generateMarkdown(pkg, reader.GetVersion())

	// Write to output file
	err = os.WriteFile(outputFile, []byte(markdown), 0644)
	if err != nil {
		log.Fatal("Failed to write markdown file:", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", siqFile, outputFile)
}

func generateMarkdown(pkg *siq.Package, version int) string {
	var sb strings.Builder

	// Process rounds
	rounds := pkg.Rounds
	if version == 4 {
		// Convert v4 rounds to v5 format for processing
		rounds = make([]siq.Round, len(pkg.RoundsV4))
		for i, roundV4 := range pkg.RoundsV4 {
			rounds[i] = convertV4RoundToV5(roundV4)
		}
	}

	for roundIndex, round := range rounds {
		sb.WriteString(fmt.Sprintf("## Round %d: %s\n\n", roundIndex+1, round.Name))

		// Process themes
		for themeIndex, theme := range round.Themes {
			sb.WriteString(fmt.Sprintf("### Theme %d: %s\n\n", themeIndex+1, theme.Name))

			// Process questions
			for questionIndex, question := range theme.Questions {
				sb.WriteString(fmt.Sprintf("#### Question %d\n\n", questionIndex+1))

				// Question content
				content := question.GetQuestionContent()
				if len(content) > 0 {
					sb.WriteString("**Content**:\n\n")
					for _, item := range content {
						sb.WriteString(fmt.Sprintf("- **%s**: %s", item.Type, item.Value))
						if item.Duration > 0 {
							sb.WriteString(fmt.Sprintf(" (duration: %d)", item.Duration))
						}
						if item.Placement != "" && item.Placement != "screen" {
							sb.WriteString(fmt.Sprintf(" (placement: %s)", item.Placement))
						}
						sb.WriteString("\n")
					}
					sb.WriteString("\n")
				}

				// Right answers
				if len(question.Right) > 0 {
					sb.WriteString("**Right Answer")
					if len(question.Right) > 1 {
						sb.WriteString("s")
					}
					sb.WriteString("**:\n\n")

					if len(question.Right) == 1 {
						sb.WriteString(question.Right[0] + "\n\n")
					} else {
						for i, answer := range question.Right {
							sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, answer))
						}
						sb.WriteString("\n")
					}
				}

				sb.WriteString("---\n\n")
			}
		}
	}

	return sb.String()
}

func convertV4RoundToV5(roundV4 siq.RoundV4) siq.Round {
	themes := make([]siq.Theme, len(roundV4.Themes))
	for i, themeV4 := range roundV4.Themes {
		questions := make([]siq.Question, len(themeV4.Questions))
		for j, questionV4 := range themeV4.Questions {
			questions[j] = siq.ConvertV4ToV5Question(questionV4)
		}
		themes[i] = siq.Theme{
			Name:      themeV4.Name,
			Info:      themeV4.Info,
			Questions: questions,
		}
	}
	return siq.Round{
		Name:   roundV4.Name,
		Info:   roundV4.Info,
		Themes: themes,
	}
}

// GetMarkdownCmd returns the markdown command
func GetMarkdownCmd() *cobra.Command {
	return markdownCmd
}
