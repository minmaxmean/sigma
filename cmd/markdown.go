package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/minmaxmean/sigma/siq"
	"github.com/spf13/cobra"
)

var (
	skipMedia bool
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

func init() {
	markdownCmd.Flags().BoolVarP(&skipMedia, "skip-media", "s", false, "Skip questions with media content (image, audio, video)")
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

// hasMediaContent checks if a question contains media content (image, audio, video)
func hasMediaContent(question *siq.Question) bool {
	content := question.GetQuestionContent()
	for _, item := range content {
		switch item.Type {
		case siq.ContentTypeImage, siq.ContentTypeAudio, siq.ContentTypeVideo, siq.ContentTypeVoice:
			return true
		case "", siq.ContentTypeMarker:
		default:
			pretty.Printf("question: %+v\n", item)
		}
	}
	return false
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
		fmt.Fprintf(&sb, "## Round %d: %s\n\n", roundIndex+1, round.Name)

		// Process themes
		for themeIndex, theme := range round.Themes {
			// Filter questions if skip-media flag is set
			var filteredQuestions []siq.Question
			if skipMedia {
				for _, question := range theme.Questions {
					if !hasMediaContent(&question) {
						filteredQuestions = append(filteredQuestions, question)
					}
				}
			} else {
				filteredQuestions = theme.Questions
			}

			// Skip themes with no questions after filtering
			if len(filteredQuestions) == 0 {
				continue
			}

			fmt.Fprintf(&sb, "### Theme %d: %s\n\n", themeIndex+1, theme.Name)

			// Process questions with original numbering
			questionNumber := 1
			for _, question := range filteredQuestions {
				fmt.Fprintf(&sb, "#### Question %d\n\n", questionNumber)

				// Question content
				content := question.GetQuestionContent()
				if len(content) > 0 {
					fmt.Fprintf(&sb, "**Content**:\n\n")
					for _, item := range content {
						fmt.Fprintf(&sb, "- %s", item.Value)
						if item.Duration > 0 {
							fmt.Fprintf(&sb, " (duration: %d)", item.Duration)
						}
						if item.Placement != "" && item.Placement != "screen" {
							fmt.Fprintf(&sb, " (placement: %s)", item.Placement)
						}
						fmt.Fprintf(&sb, "\n")
					}
					fmt.Fprintf(&sb, "\n")
				}

				// Right answers
				if len(question.Right) > 0 {
					fmt.Fprintf(&sb, "**Right Answer")
					if len(question.Right) > 1 {
						fmt.Fprintf(&sb, "s")
					}
					fmt.Fprintf(&sb, "**:\n\n")

					if len(question.Right) == 1 {
						fmt.Fprintf(&sb, "%s\n\n", question.Right[0])
					} else {
						for i, answer := range question.Right {
							fmt.Fprintf(&sb, "%d. %s\n", i+1, answer)
						}
						fmt.Fprintf(&sb, "\n")
					}
				}

				fmt.Fprintf(&sb, "---\n\n")
				questionNumber++
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
